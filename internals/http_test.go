package internals

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var _ NotionHttpClient = (*MockNotionClient)(nil) // satisfies interface

type MockNotionClient struct {
	fails     bool
	getCalls  []string
	postCalls [][]string
}

func (c *MockNotionClient) GetPage(pageId string) (string, error) {
	c.getCalls = append(c.getCalls, pageId)
	if c.fails {
		return "", errors.New("failed to get page")
	}
	return "This is a page", nil
}

func (c *MockNotionClient) PostPage(content, title, parentId string, parent ParentLiteral) (string, error) {
	c.postCalls = append(c.postCalls, []string{content, title, parentId, string(parent)})
	if c.fails {
		return "", errors.New("failed to create page")
	}
	return "some-long-id", nil
}

func NewMockNotionClient(fails bool) *MockNotionClient {
	return &MockNotionClient{
		fails: fails,
	}
}

func TestNotionClientFromDefaultsSuccess(t *testing.T) {
	t.Setenv("NOTION_API_KEY", "test-api-key")
	client, err := NewNotionClientFromDefaults()
	assert.Nil(t, err, "Error should be null")
	assert.Equal(t, client.apiKey, "test-api-key", "Incorrect API key")
	assert.Equal(t, client.notionVersion, DefaultNotionVersion, "Incorrect Notion version")
}

func TestNotionClientFromDefaultsFailure(t *testing.T) {
	t.Setenv("NOTION_API_KEY", "")
	_, err := NewNotionClientFromDefaults()
	assert.NotNil(t, err, "Error should be non-null")
	assert.Equal(t, err.Error(), "could not find NOTION_API_KEY within the current environment", "Incorrect error message")
}

func TestNotionReadSuccess(t *testing.T) {
	client := NewMockNotionClient(false)
	app := NewNotion(client)
	result, err := app.Read("some-page-id")
	assert.Nil(t, err, "Error should be null")
	assert.Equal(t, result, "This is a page", "Unexpected content")
	assert.Len(t, client.getCalls, 1, "getCalls should have a length of 1")
	assert.Equal(t, client.getCalls[0], "some-page-id", "Unexpected call argument")
}

func TestNotionReadFailure(t *testing.T) {
	client := NewMockNotionClient(true)
	app := NewNotion(client)
	_, err := app.Read("some-page-id")
	assert.NotNil(t, err, "Error should be non-null")
	assert.Equal(t, err.Error(), "failed to get page", "Unexpected error message")
	assert.Len(t, client.getCalls, 1, "getCalls should have a length of 1")
	assert.Equal(t, client.getCalls[0], "some-page-id", "Unexpected call argument")
}

func TestNotionWriteSuccess(t *testing.T) {
	client := NewMockNotionClient(false)
	app := NewNotion(client)
	result, err := app.Write("this is a page", "a page", "some-parent-id", DatabaseParentLiteral)
	assert.Nil(t, err, "Error should be null")
	assert.Equal(t, result, "some-long-id", "Unexpected content")
	assert.Len(t, client.postCalls, 1, "getCalls should have a length of 1")
	assert.Contains(t, client.postCalls[0], "this is a page", "Unexpected call argument")
	assert.Contains(t, client.postCalls[0], "a page", "Unexpected call argument")
	assert.Contains(t, client.postCalls[0], "some-parent-id", "Unexpected call argument")
	assert.Contains(t, client.postCalls[0], string(DatabaseParentLiteral), "Unexpected call argument")
}

func TestNotionWriteFailure(t *testing.T) {
	client := NewMockNotionClient(true)
	app := NewNotion(client)
	_, err := app.Write("this is a page", "a page", "some-parent-id", DatabaseParentLiteral)
	assert.NotNil(t, err, "Error should be non-null")
	assert.Equal(t, err.Error(), "failed to create page", "Unexpected error message")
	assert.Len(t, client.postCalls, 1, "getCalls should have a length of 1")
	assert.Contains(t, client.postCalls[0], "this is a page", "Unexpected call argument")
	assert.Contains(t, client.postCalls[0], "a page", "Unexpected call argument")
	assert.Contains(t, client.postCalls[0], "some-parent-id", "Unexpected call argument")
	assert.Contains(t, client.postCalls[0], string(DatabaseParentLiteral), "Unexpected call argument")
}
