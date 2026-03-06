package internals

import "time"

type PageMarkdown struct {
	Object          string   `json:"object"`
	ID              string   `json:"id"`
	Markdown        string   `json:"markdown"`
	Truncated       bool     `json:"truncated"`
	UnknownBlockIDs []string `json:"unknown_block_ids"`
}

type PostMarkdown struct {
	Markdown string `json:"markdown"`
}

type PostPage struct {
	Object         string         `json:"object"`
	ID             string         `json:"id"`
	CreatedTime    time.Time      `json:"created_time,omitempty"`
	LastEditedTime time.Time      `json:"last_edited_time,omitempty"`
	Archived       bool           `json:"archived,omitempty"`
	InTrash        bool           `json:"in_trash,omitempty"`
	IsLocked       bool           `json:"is_locked,omitempty"`
	URL            string         `json:"url,omitempty"`
	PublicURL      string         `json:"public_url,omitempty"`
	Parent         map[string]any `json:"parent,omitempty"`
	Properties     map[string]any `json:"properties,omitempty"`
	Icon           map[string]any `json:"icon,omitempty"`
	Cover          map[string]any `json:"cover,omitempty"`
	CreatedBy      User           `json:"created_by,omitempty"`
	LastEditedBy   User           `json:"last_edited_by,omitempty"`
}

type User struct {
	ID     string `json:"id,omitempty"`
	Object string `json:"object,omitempty"`
}
