package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {
	url := flag.String("url", "", "Start URL")
	depth := flag.Int("depth", 3, "Crawl depth")
	concurrency := flag.Int("concurrency", 5, "Concurrency level")
	flag.Parse()

	if *url == "" {
		log.Fatal("Missing required flag: --url")
	}

	fmt.Println("Loaded configuration:")
	fmt.Printf("  URL: %s\n", *url)
	fmt.Printf("  Depth: %d\n", *depth)
	fmt.Printf("  Concurrency: %d\n", *concurrency)
}
