package fetcher_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kvn-alcantara/gocrawler/internal/fetcher"
	"github.com/stretchr/testify/assert"
)

func TestHTTPFetcherFetchSuccess(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
			<html>
				<body>
					<a href="http://example.com/page1">Page 1</a>
					<a href="http://example.com/page2">Page 2</a>
				</body>
			</html>
		`))
	}))
	defer mockServer.Close()

	fetcher := fetcher.NewHTTPFetcher()
	links, err := fetcher.Fetch(mockServer.URL)

	assert.NoError(t, err, "expected no error")
	expectedLinks := []string{
		"http://example.com/page1",
		"http://example.com/page2",
	}
	assert.Equal(t, expectedLinks, links, "expected links to match")
}

func TestHTTPFetcherFetchNon200Response(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer mockServer.Close()

	fetcher := fetcher.NewHTTPFetcher()
	_, err := fetcher.Fetch(mockServer.URL)

	assert.Error(t, err, "expected an error")
	assert.Contains(t, err.Error(), "non-200 response", "expected error to contain 'non-200 response'")
}

func TestHTTPFetcherFetchInvalidURL(t *testing.T) {
	fetcher := fetcher.NewHTTPFetcher()
	_, err := fetcher.Fetch("http://invalid-url")

	assert.Error(t, err, "expected an error")
}

func TestHTTPFetcherFetchNoLinks(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
			<html>
				<body>
					<p>No links here!</p>
				</body>
			</html>
		`))
	}))
	defer mockServer.Close()

	fetcher := fetcher.NewHTTPFetcher()
	links, err := fetcher.Fetch(mockServer.URL)

	assert.NoError(t, err, "expected no error")
	assert.Empty(t, links, "expected no links")
}

func TestHTTPFetcherFetchBrokenHTML(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
			<html>
				<body>
					<a href="http://example.com/page1">Page 1
		`))
	}))
	defer mockServer.Close()

	fetcher := fetcher.NewHTTPFetcher()
	links, err := fetcher.Fetch(mockServer.URL)

	assert.NoError(t, err, "expected no error")
	expectedLinks := []string{
		"http://example.com/page1",
	}
	assert.Equal(t, expectedLinks, links, "expected links to match")
}

func TestHTTPFetcherTimeoutOption(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<html><body><a href="http://example.com">Link</a></body></html>`))
	}))
	defer mockServer.Close()

	shortTimeoutFetcher := fetcher.NewHTTPFetcher(fetcher.HTTPFetcherOptions{
		Timeout: 100 * time.Millisecond,
	})

	_, err := shortTimeoutFetcher.Fetch(mockServer.URL)

	assert.Error(t, err, "expected timeout error")
	assert.Contains(t, err.Error(), "Client.Timeout exceeded", "error should indicate timeout")
}
