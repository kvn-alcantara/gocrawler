package crawler_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/kvn-alcantara/gocrawler/internal/crawler"
	"github.com/stretchr/testify/assert"
)

type MockFetcher struct {
	links map[string][]string
	error error
	mu    sync.Mutex
}

func (m *MockFetcher) Fetch(url string) ([]string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.links[url], nil
}

func TestCrawl(t *testing.T) {
	mockFetcher := &MockFetcher{
		links: map[string][]string{
			"http://example.com": {
				"http://example.com/page1",
				"http://example.com/page2",
			},
			"http://example.com/page1": {
				"http://example.com/page3",
			},
			"http://example.com/page2": {},
			"http://example.com/page3": {},
		},
	}

	cfg := crawler.Config{
		StartURL:    "http://example.com",
		MaxDepth:    2,
		Concurrency: 2,
		Fetcher:     mockFetcher,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	visited := crawler.Crawl(ctx, cfg)

	expectedVisited := map[string]bool{
		"http://example.com":       true,
		"http://example.com/page1": true,
		"http://example.com/page2": true,
		"http://example.com/page3": true,
	}

	assert.Equal(t, expectedVisited, visited, "visited URLs should match expected")
}

func TestCrawlWithMaxDepth(t *testing.T) {
	mockFetcher := &MockFetcher{
		links: map[string][]string{
			"http://example.com": {
				"http://example.com/page1",
				"http://example.com/page2",
			},
			"http://example.com/page1": {
				"http://example.com/page3",
			},
			"http://example.com/page2": {},
			"http://example.com/page3": {},
		},
	}

	cfg := crawler.Config{
		StartURL:    "http://example.com",
		MaxDepth:    1,
		Concurrency: 2,
		Fetcher:     mockFetcher,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	visited := crawler.Crawl(ctx, cfg)

	expectedVisited := map[string]bool{
		"http://example.com":       true,
		"http://example.com/page1": true,
		"http://example.com/page2": true,
	}

	assert.Equal(t, expectedVisited, visited, "visited URLs should match expected")
}
