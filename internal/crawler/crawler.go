package crawler

import (
	"context"
	"fmt"
	"sync"

	"github.com/kvn-alcantara/gocrawler/internal/fetcher"
)

// Config holds the configuration for the web crawler.
type Config struct {
	StartURL    string
	MaxDepth    int
	Concurrency int
	Fetcher     fetcher.PageFetcher
}

// Crawl initializes the visited map and starts the recursive crawling.
func Crawl(ctx context.Context, cfg Config) map[string]bool {
	visited := make(map[string]bool)
	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, cfg.Concurrency)

	wg.Add(1)
	go func() {
		defer wg.Done()
		crawl(ctx, cfg.StartURL, 0, cfg, visited, &mu, &wg, sem)
	}()

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		fmt.Println("Crawling finished.")
	case <-ctx.Done():
		fmt.Println("Crawling canceled.")
	}

	return visited
}

// crawl performs the recursive crawling of URLs.
func crawl(ctx context.Context, url string, depth int, cfg Config, visited map[string]bool, mu *sync.Mutex, wg *sync.WaitGroup, sem chan struct{}) {
	if depth > cfg.MaxDepth {
		return
	}

	mu.Lock()
	if visited[url] {
		mu.Unlock()
		return
	}
	visited[url] = true
	fmt.Println("Crawling:", url, "at depth", depth)
	mu.Unlock()

	wg.Add(1)
	sem <- struct{}{}

	go func() {
		defer func() {
			<-sem
			wg.Done()
		}()

		select {
		case <-ctx.Done():
			fmt.Println("Crawl canceled before fetching:", url)
			return
		default:
		}

		links, err := cfg.Fetcher.Fetch(url)
		if err != nil {
			fmt.Println("Error fetching", url, ":", err)
			return
		}

		for _, link := range links {
			select {
			case <-ctx.Done():
				fmt.Println("Crawl canceled while processing links from:", url)
				return
			default:
			}

			mu.Lock()
			if !visited[link] {
				mu.Unlock()
				wg.Add(1)
				go func() {
					defer wg.Done()
					crawl(ctx, link, depth+1, cfg, visited, mu, wg, sem)
				}()
			} else {
				mu.Unlock()
			}
		}
	}()
}
