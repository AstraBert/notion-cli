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

const DefaultNotionVersion string = "2025-09-03"

var _ NotionHttpClient = (*NotionClient)(nil) // satisfies interface

type NotionHttpClient interface {
	GetPage(string) (string, error)
	PostPage(string) (string, error)
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
	if !ok {
		return nil, errors.New("could not find NOTION_API_KEY within the current environment")
	}
	return &NotionClient{
		apiKey:        apiKey,
		notionVersion: DefaultNotionVersion,
	}, nil
}

func (n *NotionClient) GetPage(pageId string) (string, error) {
	client := &http.Client{Timeout: time.Duration(60) * time.Second}
	url := fmt.Sprintf("https://api.notion.com/v1/pages/%s/markdown", pageId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Notion-Version", n.notionVersion)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer: %s", n.apiKey))

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	if res.StatusCode > 299 || res.StatusCode < 200 {
		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return "", err
		}
		detail := string(body)
		return "", fmt.Errorf("Response returned a status code of %d: %s", res.StatusCode, detail)
	}
	defer res.Body.Close()
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

func (n *NotionClient) PostPage(markdownContent string) (string, error) {
	client := &http.Client{Timeout: time.Duration(60) * time.Second}
	url := "https://api.notion.com/v1/pages"
	reqBodyJson := PostMarkdown{Markdown: markdownContent}
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
	req.Header.Add("Authorization", fmt.Sprintf("Bearer: %s", n.apiKey))
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	if res.StatusCode > 299 || res.StatusCode < 200 {
		defer res.Body.Close()
		respBody, err := io.ReadAll(res.Body)
		if err != nil {
			return "", err
		}
		detail := string(respBody)
		return "", fmt.Errorf("Response returned a status code of %d: %s", res.StatusCode, detail)
	}
	defer res.Body.Close()
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

type Notion struct {
	client NotionHttpClient
}

func NewNotion(client NotionHttpClient) *Notion {
	return &Notion{
		client: client,
	}
}

func (app *Notion) Read(pageId string) (string, error) {
	return app.client.GetPage(pageId)
}

func (app *Notion) Write(content string) (string, error) {
	return app.client.PostPage(content)
}
