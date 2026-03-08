package internals

import (
	"fmt"
	"math/rand/v2"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNotionReadSuccessIntegration(t *testing.T) {
	apiKey, okApi := os.LookupEnv("TEST_NOTION_API_KEY")
	if !okApi {
		t.Skip("NOTION_API_KEY not available")
	}
	pageId, okPage := os.LookupEnv("NOTION_PAGE_ID")
	if !okPage {
		t.Skip("NOTION_PAGE_ID not available")
	}
	client := NewNotionClient(apiKey, DefaultNotionVersion)
	app := NewNotion(client)
	content, err := app.Read(pageId, MaxRetries, DefaultRetryTime)
	assert.Nil(t, err, "Error should be null")
	assert.Greater(t, len(content), 0, "Content should be non-empty")
}

func TestNotionReadFailureIntegration(t *testing.T) {
	apiKey, okApi := os.LookupEnv("TEST_NOTION_API_KEY")
	if !okApi {
		t.Skip("NOTION_API_KEY not available")
	}
	client := NewNotionClient(apiKey, DefaultNotionVersion)
	app := NewNotion(client)
	_, err := app.Read("invalid-page-id", MaxRetries, DefaultRetryTime)
	assert.NotNil(t, err, "Error should be non-null")
	assert.Contains(t, err.Error(), "response returned a status code of 400:", "Unexpected error message")
}

func TestNotionWriteSuccessIntegration(t *testing.T) {
	apiKey, okApi := os.LookupEnv("TEST_NOTION_API_KEY")
	if !okApi {
		t.Skip("NOTION_API_KEY not available")
	}
	pageId, okPage := os.LookupEnv("NOTION_PAGE_ID")
	if !okPage {
		t.Skip("NOTION_PAGE_ID not available")
	}
	client := NewNotionClient(apiKey, DefaultNotionVersion)
	app := NewNotion(client)
	randNum := rand.IntN(1000)
	content, err := app.Write("some content", fmt.Sprintf("page-%d", randNum), pageId, PageParentLiteral, MaxRetries, DefaultRetryTime)
	assert.Nil(t, err, "Error should be null")
	assert.Equal(t, len(content), 36, "Length of a UUIDv4 should be 36")
	assert.Equal(t, strings.Count(content, "-"), 4, "UUIDv4 should contain 4 dashes")
}

func TestNotionWriteFailureIntegration(t *testing.T) {
	apiKey, okApi := os.LookupEnv("TEST_NOTION_API_KEY")
	if !okApi {
		t.Skip("NOTION_API_KEY not available")
	}
	client := NewNotionClient(apiKey, DefaultNotionVersion)
	app := NewNotion(client)
	randNum := rand.IntN(1000)
	_, err := app.Write("some content", fmt.Sprintf("page-%d", randNum), "invalid-uuid", PageParentLiteral, MaxRetries, DefaultRetryTime)
	assert.NotNil(t, err, "Error should be non-null")
	assert.Contains(t, err.Error(), "response returned a status code of 400:", "Unexpected error message")
}

func TestNotionAppendSuccessIntegration(t *testing.T) {
	apiKey, okApi := os.LookupEnv("TEST_NOTION_API_KEY")
	if !okApi {
		t.Skip("NOTION_API_KEY not available")
	}
	pageId, okPage := os.LookupEnv("NOTION_PAGE_ID")
	if !okPage {
		t.Skip("NOTION_PAGE_ID not available")
	}
	client := NewNotionClient(apiKey, DefaultNotionVersion)
	app := NewNotion(client)
	randNum := rand.IntN(1000)
	returnedId, err := app.Append(pageId, fmt.Sprintf("Thank you %d", randNum), MaxRetries, DefaultRetryTime)
	assert.Nil(t, err, "Error should be null")
	assert.Equal(t, returnedId, pageId, "Original page ID and returned ID should match")
}

func TestNotionAppendFailureIntegration(t *testing.T) {
	apiKey, okApi := os.LookupEnv("TEST_NOTION_API_KEY")
	if !okApi {
		t.Skip("NOTION_API_KEY not available")
	}
	client := NewNotionClient(apiKey, DefaultNotionVersion)
	app := NewNotion(client)
	randNum := rand.IntN(1000)
	_, err := app.Append("invalid-uuid", fmt.Sprintf("Thank you %d", randNum), MaxRetries, DefaultRetryTime)
	assert.NotNil(t, err, "Error should be non-null")
	assert.Contains(t, err.Error(), "response returned a status code of 400:", "Unexpected error message")
}
