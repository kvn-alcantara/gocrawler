# GoCrawler

## Buil
```bash
go build ./cmd/webcrawler
```

## Example
```bash
./webcrawler --url https://github.com --depth 2 --concurrency 5
```

## Flags
`--url`: The URL to start crawling from.
`--depth`: The maximum depth to crawl. Default is 3.
`--concurrency`: The number of concurrent requests to make. Default is 5.

## Test
```bash
go test ./...
```
