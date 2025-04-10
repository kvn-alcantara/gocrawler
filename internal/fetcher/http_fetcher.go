package fetcher

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/html"
)

// HTTPFetcher implements the PageFetcher interface for fetching web pages.
type HTTPFetcher struct {
	client *http.Client
}

// HTTPFetcherOptions defines configuration options for the HTTPFetcher
type HTTPFetcherOptions struct {
	Timeout time.Duration
}

// NewHTTPFetcher creates a new instance of HTTPFetcher.
func NewHTTPFetcher(options ...HTTPFetcherOptions) *HTTPFetcher {
	var opts HTTPFetcherOptions
	if len(options) > 0 {
		opts = options[0]
	}

	timeout := opts.Timeout
	if timeout == 0 {
		timeout = 10 * time.Second
	}

	return &HTTPFetcher{
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

// Fetch fetches the content of the given URL and returns a list of links found in the HTML.
func (f *HTTPFetcher) Fetch(url string) ([]string, error) {
	res, err := f.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, errors.New("non-200 response: " + res.Status)
	}

	tokenizer := html.NewTokenizer(res.Body)
	var links []string

	for {
		tt := tokenizer.Next()
		if tt == html.ErrorToken {
			break
		}

		token := tokenizer.Token()
		if token.Data == "a" {
			for _, attr := range token.Attr {
				if attr.Key == "href" && strings.HasPrefix(attr.Val, "http") {
					links = append(links, attr.Val)
				}
			}
		}
	}
	return links, nil
}

var _ PageFetcher = &HTTPFetcher{}
