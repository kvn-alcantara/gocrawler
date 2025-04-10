package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kvn-alcantara/gocrawler/internal/crawler"
	"github.com/kvn-alcantara/gocrawler/internal/fetcher"
)

func main() {
	url := flag.String("url", "", "Start URL")
	depth := flag.Int("depth", 3, "Crawl depth")
	concurrency := flag.Int("concurrency", 5, "Concurrency level")
	timeout := flag.Int("timeout", 10, "HTTP request timeout in seconds")
	flag.Parse()

	if *url == "" {
		log.Fatal("Missing required flag: --url")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChan
		fmt.Println("Interrupt received, shutting down...")
		cancel()
	}()

	fetcherOptions := fetcher.HTTPFetcherOptions{
		Timeout: time.Duration(*timeout) * time.Second,
	}

	config := crawler.Config{
		StartURL:    *url,
		MaxDepth:    *depth,
		Concurrency: *concurrency,
		Fetcher:     fetcher.NewHTTPFetcher(fetcherOptions),
	}

	visited := crawler.Crawl(ctx, config)

	log.Println("Visited URLs:")
	for url := range visited {
		log.Println(url)
	}
}
