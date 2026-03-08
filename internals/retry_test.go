package internals

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func CreateHandler() http.Handler {
	mux := http.NewServeMux()
	unstableCalls := 0

	mux.HandleFunc("GET /stable", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("<h1>Hello</h1>"))
	})

	mux.HandleFunc("GET /unstable", func(w http.ResponseWriter, r *http.Request) {
		unstableCalls += 1
		if unstableCalls <= 2 {
			w.Header().Set("Retry-After", "1")
			http.Error(
				w,
				"Too Many Requests",
				http.StatusTooManyRequests,
			)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("<h1>Hello</h1>"))
	})

	mux.HandleFunc("GET /error", func(w http.ResponseWriter, r *http.Request) {
		http.Error(
			w,
			"Server replied with an error",
			http.StatusInternalServerError,
		)
	})

	mux.HandleFunc("GET /unretriable", func(w http.ResponseWriter, r *http.Request) {
		http.Error(
			w,
			"Bad request",
			http.StatusBadRequest,
		)
	})

	return mux
}

func TestRetryStable(t *testing.T) {
	server := httptest.NewServer(CreateHandler())
	defer server.Close()
	client := server.Client()
	req, _ := http.NewRequest("GET", server.URL+"/stable", nil)
	startTime := time.Now()
	resp, err := RequestWithRetries(client, req, MaxRetries, DefaultRetryTime)
	assert.Nil(t, err, "Error should be null")
	assert.Equal(t, resp.StatusCode, http.StatusOK, "Should return an OK response")
	assert.Less(t, time.Since(startTime), time.Duration(1)*time.Second, "Should take less than one second")
	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("An error occurred while reading the body of the response: %s\n", err.Error())
	}
	content := string(body)
	assert.Equal(t, content, "<h1>Hello</h1>")
}

func TestRetryUnstable(t *testing.T) {
	server := httptest.NewServer(CreateHandler())
	defer server.Close()
	client := server.Client()
	req, _ := http.NewRequest("GET", server.URL+"/unstable", nil)
	startTime := time.Now()
	resp, err := RequestWithRetries(client, req, MaxRetries, DefaultRetryTime)
	assert.Nil(t, err, "Error should be null")
	assert.Equal(t, resp.StatusCode, http.StatusOK, "Should return an OK response")
	assert.GreaterOrEqual(t, time.Since(startTime), time.Duration(2)*time.Second, "Should take two seconds or more")
	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("An error occurred while reading the body of the response: %s\n", err.Error())
	}
	content := string(body)
	assert.Equal(t, content, "<h1>Hello</h1>")
}

func TestRetryError(t *testing.T) {
	server := httptest.NewServer(CreateHandler())
	defer server.Close()
	client := server.Client()
	req, _ := http.NewRequest("GET", server.URL+"/error", nil)
	startTime := time.Now()
	_, err := RequestWithRetries(client, req, MaxRetries, DefaultRetryTime)
	assert.NotNil(t, err, "Error should be not-null")
	assert.Equal(t, err.Error(), "exceeded maximum number of retries: 3", "Unexpected error message")
	assert.GreaterOrEqual(t, time.Since(startTime), time.Duration(3)*time.Second, "Should take 3 seconds or more")
}

func TestRetryUnretriable(t *testing.T) {
	server := httptest.NewServer(CreateHandler())
	defer server.Close()
	client := server.Client()
	req, _ := http.NewRequest("GET", server.URL+"/unretriable", nil)
	startTime := time.Now()
	resp, err := RequestWithRetries(client, req, MaxRetries, DefaultRetryTime)
	assert.Nil(t, err, "Error should be null")
	assert.Equal(t, resp.StatusCode, http.StatusBadRequest, "Should return a 400 response")
	assert.Less(t, time.Since(startTime), time.Duration(1)*time.Second, "Should take less than 1 second")
	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("An error occurred while reading the body of the response: %s\n", err.Error())
	}
	content := string(body)
	assert.Equal(t, content, "Bad request\n")
}
