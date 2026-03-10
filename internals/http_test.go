package internals

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var _ NotionHttpClient = (*MockNotionClient)(nil) // satisfies interface

type MockNotionClient struct {
	fails         bool
	getCalls      []string
	postCalls     [][]string
	patchCalls    [][]string
	retryArgs     []int
	searchQueries []string
}

func (c *MockNotionClient) GetPage(pageId string, maxRetries, retryTime int) (string, error) {
	c.getCalls = append(c.getCalls, pageId)
	c.retryArgs = append(c.retryArgs, maxRetries, retryTime)
	if c.fails {
		return "", errors.New("failed to get page")
	}
	return "This is a page", nil
}

func (c *MockNotionClient) PostPage(content, title, parentId string, parent ParentLiteral, maxRetries, retryTime int) (string, error) {
	c.postCalls = append(c.postCalls, []string{content, title, parentId, string(parent)})
	c.retryArgs = append(c.retryArgs, maxRetries, retryTime)
	if c.fails {
		return "", errors.New("failed to create page")
	}
	return "some-long-id", nil
}

func (c *MockNotionClient) PatchPage(pageId string, content string, maxRetries, retryTime int) (string, error) {
	c.patchCalls = append(c.patchCalls, []string{pageId, content})
	c.retryArgs = append(c.retryArgs, maxRetries, retryTime)
	if c.fails {
		return "", errors.New("failed to patch page")
	}
	return "patched-long-id", nil
}

func (c *MockNotionClient) SearchPages(query, startCursor string, sortStrategy SortStrategyLiteral, pageSize, maxRetries, retryInterval int) ([]string, error) {
	c.searchQueries = append(c.searchQueries, query)
	if c.fails {
		return nil, errors.New("failed to search pages")
	}
	return []string{"hello-world", "bye-moon"}, nil
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
	result, err := app.Read("some-page-id", 3, 1)
	assert.Nil(t, err, "Error should be null")
	assert.Equal(t, result, "This is a page", "Unexpected content")
	assert.Len(t, client.getCalls, 1, "getCalls should have a length of 1")
	assert.Equal(t, client.getCalls[0], "some-page-id", "Unexpected call argument")
	assert.Len(t, client.retryArgs, 2, "retryArgs should have a length of 2")
	assert.Equal(t, client.retryArgs[0], 3, "maxRetries should be 3")
	assert.Equal(t, client.retryArgs[1], 1, "retryInterval should be 1")
}

func TestNotionReadFailure(t *testing.T) {
	client := NewMockNotionClient(true)
	app := NewNotion(client)
	_, err := app.Read("some-page-id", 3, 1)
	assert.NotNil(t, err, "Error should be non-null")
	assert.Equal(t, err.Error(), "failed to get page", "Unexpected error message")
	assert.Len(t, client.getCalls, 1, "getCalls should have a length of 1")
	assert.Equal(t, client.getCalls[0], "some-page-id", "Unexpected call argument")
	assert.Len(t, client.retryArgs, 2, "retryArgs should have a length of 2")
	assert.Equal(t, client.retryArgs[0], 3, "maxRetries should be 3")
	assert.Equal(t, client.retryArgs[1], 1, "retryInterval should be 1")
}

func TestNotionWriteSuccess(t *testing.T) {
	client := NewMockNotionClient(false)
	app := NewNotion(client)
	result, err := app.Write("this is a page", "a page", "some-parent-id", DatabaseParentLiteral, 3, 1)
	assert.Nil(t, err, "Error should be null")
	assert.Equal(t, result, "some-long-id", "Unexpected content")
	assert.Len(t, client.postCalls, 1, "postCalls should have a length of 1")
	assert.Contains(t, client.postCalls[0], "this is a page", "Unexpected call argument")
	assert.Contains(t, client.postCalls[0], "a page", "Unexpected call argument")
	assert.Contains(t, client.postCalls[0], "some-parent-id", "Unexpected call argument")
	assert.Contains(t, client.postCalls[0], string(DatabaseParentLiteral), "Unexpected call argument")
	assert.Len(t, client.retryArgs, 2, "retryArgs should have a length of 2")
	assert.Equal(t, client.retryArgs[0], 3, "maxRetries should be 3")
	assert.Equal(t, client.retryArgs[1], 1, "retryInterval should be 1")
}

func TestNotionWriteFailure(t *testing.T) {
	client := NewMockNotionClient(true)
	app := NewNotion(client)
	_, err := app.Write("this is a page", "a page", "some-parent-id", DatabaseParentLiteral, 3, 1)
	assert.NotNil(t, err, "Error should be non-null")
	assert.Equal(t, err.Error(), "failed to create page", "Unexpected error message")
	assert.Len(t, client.postCalls, 1, "postCalls should have a length of 1")
	assert.Contains(t, client.postCalls[0], "this is a page", "Unexpected call argument")
	assert.Contains(t, client.postCalls[0], "a page", "Unexpected call argument")
	assert.Contains(t, client.postCalls[0], "some-parent-id", "Unexpected call argument")
	assert.Contains(t, client.postCalls[0], string(DatabaseParentLiteral), "Unexpected call argument")
	assert.Len(t, client.retryArgs, 2, "retryArgs should have a length of 2")
	assert.Equal(t, client.retryArgs[0], 3, "maxRetries should be 3")
	assert.Equal(t, client.retryArgs[1], 1, "retryInterval should be 1")
}

func TestNotionAppendSuccess(t *testing.T) {
	client := NewMockNotionClient(false)
	app := NewNotion(client)
	result, err := app.Append("page-id", "appended content", 3, 1)
	assert.Nil(t, err, "Error should be null")
	assert.Equal(t, result, "patched-long-id", "Unexpected content")
	assert.Len(t, client.patchCalls, 1, "patchCalls should have a length of 1")
	assert.Contains(t, client.patchCalls[0], "page-id", "Unexpected call argument")
	assert.Contains(t, client.patchCalls[0], "appended content", "Unexpected call argument")
	assert.Len(t, client.retryArgs, 2, "retryArgs should have a length of 2")
	assert.Equal(t, client.retryArgs[0], 3, "maxRetries should be 3")
	assert.Equal(t, client.retryArgs[1], 1, "retryInterval should be 1")
}

func TestNotionAppendFailure(t *testing.T) {
	client := NewMockNotionClient(true)
	app := NewNotion(client)
	_, err := app.Append("page-id", "appended content", 3, 1)
	assert.NotNil(t, err, "Error should be non-null")
	assert.Equal(t, err.Error(), "failed to patch page", "Unexpected error message")
	assert.Len(t, client.patchCalls, 1, "patchCalls should have a length of 1")
	assert.Contains(t, client.patchCalls[0], "page-id", "Unexpected call argument")
	assert.Contains(t, client.patchCalls[0], "appended content", "Unexpected call argument")
	assert.Len(t, client.retryArgs, 2, "retryArgs should have a length of 2")
	assert.Equal(t, client.retryArgs[0], 3, "maxRetries should be 3")
	assert.Equal(t, client.retryArgs[1], 1, "retryInterval should be 1")
}

func TestNotionSearchSuccess(t *testing.T) {
	client := NewMockNotionClient(false)
	app := NewNotion(client)
	result, err := app.Search("hello", "", AscendingSortStrategy, -1, 3, 1)
	assert.Nil(t, err, "Error should be null")
	assert.Len(t, result, 2, "Unexpected result length")
	assert.Contains(t, result, "hello-world", "Unexpected result element")
	assert.Contains(t, result, "bye-moon", "Unexpected result element")
	assert.Len(t, client.searchQueries, 1, "searchQueries should have a length of 1")
	assert.Equal(t, client.searchQueries[0], "hello", "Unexpected call argument")
}

func TestNotionSearchFailure(t *testing.T) {
	client := NewMockNotionClient(true)
	app := NewNotion(client)
	_, err := app.Search("hello", "", AscendingSortStrategy, -1, 3, 1)
	assert.NotNil(t, err, "Error should be non-null")
	assert.Equal(t, err.Error(), "failed to search pages", "Unexpected error message")
	assert.Len(t, client.searchQueries, 1, "searchQueries should have a length of 1")
	assert.Equal(t, client.searchQueries[0], "hello", "Unexpected call argument")
}
