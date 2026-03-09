package internals

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type ParentLiteral string
type SortStrategyLiteral string

const DatabaseParentLiteral ParentLiteral = "database"
const PageParentLiteral ParentLiteral = "page"
const AscendingSortStrategy SortStrategyLiteral = "ascending"
const DescendingSortStrategy SortStrategyLiteral = "descending"
const DefaultSortBy string = "last_edited_time"
const DefaultFilterValue string = "page"
const DefaultFilterProperty string = "object"
const DefaultNotionVersion string = "2025-09-03"

var _ NotionHttpClient = (*NotionClient)(nil) // satisfies interface

type NotionHttpClient interface {
	GetPage(string, int, int) (string, error)
	PostPage(string, string, string, ParentLiteral, int, int) (string, error)
	PatchPage(string, string, int, int) (string, error)
	SearchPages(string, string, SortStrategyLiteral, int, int, int) ([]string, error)
}

type NotionClient struct {
	apiKey        string
	notionVersion string
}

func NewNotionClient(apiKey string, notionVersion string) *NotionClient {
	return &NotionClient{
		apiKey:        apiKey,
		notionVersion: notionVersion,
	}
}

func NewNotionClientFromDefaults() (*NotionClient, error) {
	apiKey, ok := os.LookupEnv("NOTION_API_KEY")
	if !ok || apiKey == "" {
		return nil, errors.New("could not find NOTION_API_KEY within the current environment")
	}
	return &NotionClient{
		apiKey:        apiKey,
		notionVersion: DefaultNotionVersion,
	}, nil
}

func (n *NotionClient) GetPage(pageId string, maxRetries, retryInterval int) (string, error) {
	client := &http.Client{Timeout: time.Duration(60) * time.Second}
	url := fmt.Sprintf("https://api.notion.com/v1/pages/%s/markdown", pageId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Notion-Version", n.notionVersion)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", n.apiKey))

	res, err := RequestWithRetries(client, req, maxRetries, retryInterval)
	if err != nil {
		return "", err
	}
	if res.StatusCode > 299 || res.StatusCode < 200 {
		defer func() { _ = res.Body.Close() }()
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return "", err
		}
		detail := string(body)
		return "", fmt.Errorf("response returned a status code of %d: %s", res.StatusCode, detail)
	}
	defer func() { _ = res.Body.Close() }()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	var page PageMarkdown
	err = json.Unmarshal(body, &page)
	if err != nil {
		return "", err
	}
	return page.Markdown, nil
}

func (n *NotionClient) PostPage(markdownContent, title, parentId string, parentType ParentLiteral, maxRetries, retryInterval int) (string, error) {
	client := &http.Client{Timeout: time.Duration(60) * time.Second}
	url := "https://api.notion.com/v1/pages"
	var parent ParentType
	switch parentType {
	case DatabaseParentLiteral:
		parent = &DatabaseParent{Type: "database_id", DatabaseId: parentId}
	case PageParentLiteral:
		parent = &PageParent{Type: "page_id", PageId: parentId}
	}
	reqBodyJson := PostMarkdown{
		Markdown: markdownContent,
		Parent:   parent,
		Properties: PageProperties{
			Title: TitleProperty{
				Title: []RichTextItem{
					{
						Type: "text",
						Text: RichTextBody{
							Content: title,
						},
					},
				},
			},
		},
	}
	bodyData, err := json.Marshal(reqBodyJson)
	if err != nil {
		return "", err
	}
	body := bytes.NewReader(bodyData)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return "", err
	}
	req.Header.Add("Notion-Version", n.notionVersion)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", n.apiKey))
	req.Header.Add("Content-Type", "application/json")

	res, err := RequestWithRetries(client, req, maxRetries, retryInterval)
	if err != nil {
		return "", err
	}
	if res.StatusCode > 299 || res.StatusCode < 200 {
		defer func() { _ = res.Body.Close() }()
		respBody, err := io.ReadAll(res.Body)
		if err != nil {
			return "", err
		}
		detail := string(respBody)
		return "", fmt.Errorf("response returned a status code of %d: %s", res.StatusCode, detail)
	}
	defer func() { _ = res.Body.Close() }()
	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	var page PostPage
	err = json.Unmarshal(respBody, &page)
	if err != nil {
		return "", err
	}
	return page.ID, nil
}

func (n *NotionClient) PatchPage(pageId string, content string, maxRetries, retryInterval int) (string, error) {
	client := &http.Client{Timeout: time.Duration(60) * time.Second}
	url := fmt.Sprintf("https://api.notion.com/v1/pages/%s/markdown", pageId)
	reqBodyJson := PatchMarkdown{Type: "insert_content", InsertContent: InsertContent{Content: content}}
	bodyData, err := json.Marshal(reqBodyJson)
	if err != nil {
		return "", err
	}
	body := bytes.NewReader(bodyData)
	req, err := http.NewRequest("PATCH", url, body)
	if err != nil {
		return "", err
	}
	req.Header.Add("Notion-Version", n.notionVersion)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", n.apiKey))
	req.Header.Add("Content-Type", "application/json")

	res, err := RequestWithRetries(client, req, maxRetries, retryInterval)
	if err != nil {
		return "", err
	}
	if res.StatusCode > 299 || res.StatusCode < 200 {
		defer func() { _ = res.Body.Close() }()
		respBody, err := io.ReadAll(res.Body)
		if err != nil {
			return "", err
		}
		detail := string(respBody)
		return "", fmt.Errorf("response returned a status code of %d: %s", res.StatusCode, detail)
	}
	defer func() { _ = res.Body.Close() }()
	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	var page PatchedPage
	err = json.Unmarshal(respBody, &page)
	if err != nil {
		return "", err
	}
	return page.ID, nil
}

func (n *NotionClient) SearchPages(query, startCursor string, sortStrategy SortStrategyLiteral, pageSize, maxRetries, retryInterval int) ([]string, error) {
	client := &http.Client{Timeout: time.Duration(60) * time.Second}
	url := "https://api.notion.com/v1/search"
	reqBodyJson := SearchPagesRequest{
		Query: query,
		Sort: SearchSortBy{
			Timestamp: DefaultSortBy,
			Direction: sortStrategy,
		},
		Filter: SearchFilter{
			Property: DefaultFilterProperty,
			Value:    DefaultFilterValue,
		},
	}
	if pageSize != -1 {
		reqBodyJson.PageSize = pageSize
	}
	if startCursor != "" {
		reqBodyJson.StartCursor = startCursor
	}
	bodyData, err := json.Marshal(reqBodyJson)
	if err != nil {
		return nil, err
	}
	body := bytes.NewReader(bodyData)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Notion-Version", n.notionVersion)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", n.apiKey))
	req.Header.Add("Content-Type", "application/json")

	res, err := RequestWithRetries(client, req, maxRetries, retryInterval)
	if err != nil {
		return nil, err
	}
	if res.StatusCode > 299 || res.StatusCode < 200 {
		defer func() { _ = res.Body.Close() }()
		respBody, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		detail := string(respBody)
		return nil, fmt.Errorf("response returned a status code of %d: %s", res.StatusCode, detail)
	}
	defer func() { _ = res.Body.Close() }()
	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var results SearchPagesResponse
	err = json.Unmarshal(respBody, &results)
	if err != nil {
		return nil, err
	}
	ids := []string{}
	for i := range results.Results {
		ids = append(ids, results.Results[i].ID)
	}
	if results.HasMore && results.NextCursor != nil {
		newResults, err := n.SearchPages(query, *results.NextCursor, sortStrategy, pageSize, maxRetries, retryInterval)
		if err != nil {
			return ids, err
		}
		ids = append(ids, newResults...)
	}
	return ids, nil
}

type Notion struct {
	client NotionHttpClient
}

func NewNotion(client NotionHttpClient) *Notion {
	return &Notion{
		client: client,
	}
}

func (app *Notion) Read(pageId string, maxRetries, retryInterval int) (string, error) {
	return app.client.GetPage(pageId, maxRetries, retryInterval)
}

func (app *Notion) Write(content, title, parentId string, parentType ParentLiteral, maxRetries, retryInterval int) (string, error) {
	return app.client.PostPage(content, title, parentId, parentType, maxRetries, retryInterval)
}

func (app *Notion) Append(pageId string, content string, maxRetries, retryInterval int) (string, error) {
	return app.client.PatchPage(pageId, content, maxRetries, retryInterval)
}

func (app *Notion) Search(query, startCursor string, sortStrategy SortStrategyLiteral, pageSize, maxRetries, retryInterval int) ([]string, error) {
	return app.client.SearchPages(query, startCursor, sortStrategy, pageSize, maxRetries, retryInterval)
}
