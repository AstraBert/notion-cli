package internals

import "time"

var _ ParentType = (*PageParent)(nil)     // satisfies interface
var _ ParentType = (*DatabaseParent)(nil) // satisfies interface

type PageMarkdown struct {
	Object          string   `json:"object"`
	ID              string   `json:"id"`
	Markdown        string   `json:"markdown"`
	Truncated       bool     `json:"truncated"`
	UnknownBlockIDs []string `json:"unknown_block_ids"`
}

type PostMarkdown struct {
	Markdown   string         `json:"markdown"`
	Parent     ParentType     `json:"parent"`
	Properties PageProperties `json:"properties"`
}

type ParentType interface {
	GetId() string
}

type PageParent struct {
	Type   string `json:"type"`
	PageId string `json:"page_id"`
}

func (p *PageParent) GetId() string {
	return p.PageId
}

type DatabaseParent struct {
	Type       string `json:"type"`
	DatabaseId string `json:"database_id"`
}

func (d *DatabaseParent) GetId() string {
	return d.DatabaseId
}

type PageProperties struct {
	Title TitleProperty `json:"title"`
}

type TitleProperty struct {
	Title []RichTextItem `json:"title"`
}

type RichTextItem struct {
	Type string       `json:"type"`
	Text RichTextBody `json:"text"`
}

type RichTextBody struct {
	Content string `json:"content"`
}

type PostPage struct {
	Object         string         `json:"object"`
	ID             string         `json:"id"`
	CreatedTime    time.Time      `json:"created_time"`
	LastEditedTime time.Time      `json:"last_edited_time"`
	Archived       bool           `json:"archived,omitempty"`
	InTrash        bool           `json:"in_trash,omitempty"`
	IsLocked       bool           `json:"is_locked,omitempty"`
	URL            string         `json:"url,omitempty"`
	PublicURL      string         `json:"public_url,omitempty"`
	Parent         map[string]any `json:"parent,omitempty"`
	Properties     map[string]any `json:"properties,omitempty"`
	Icon           map[string]any `json:"icon,omitempty"`
	Cover          map[string]any `json:"cover,omitempty"`
	CreatedBy      User           `json:"created_by"`
	LastEditedBy   User           `json:"last_edited_by"`
}

type User struct {
	ID     string `json:"id,omitempty"`
	Object string `json:"object,omitempty"`
}

type PatchMarkdown struct {
	Type          string        `json:"type"`
	InsertContent InsertContent `json:"insert_content"`
}

type InsertContent struct {
	Content string `json:"content"`
}

type PatchedPage struct {
	Object          string   `json:"object"`
	ID              string   `json:"id"`
	Markdown        string   `json:"markdown"`
	Truncated       bool     `json:"truncated"`
	UnknownBlockIDs []string `json:"unknown_block_ids"`
}

type SearchPagesRequest struct {
	Query       string       `json:"query"`
	Sort        SearchSortBy `json:"sort"`
	Filter      SearchFilter `json:"filter"`
	PageSize    int          `json:"page_size,omitempty"`
	StartCursor string       `json:"start_cursor,omitempty"`
}

type SearchSortBy struct {
	Timestamp string              `json:"timestamp"`
	Direction SortStrategyLiteral `json:"direction"`
}

type SearchFilter struct {
	Property string `json:"property"`
	Value    string `json:"value"`
}

type SearchPagesResponse struct {
	Type             string         `json:"type"`
	PageOrDataSource map[string]any `json:"page_or_data_source"`
	Object           string         `json:"object"`
	NextCursor       *string        `json:"next_cursor"`
	HasMore          bool           `json:"has_more"`
	Results          []PostPage     `json:"results"`
}
