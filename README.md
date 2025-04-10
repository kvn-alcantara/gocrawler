# GoCrawler

## Build
```bash
go build ./cmd/webcrawler
```

## Example
```bash
./webcrawler --url https://github.com --depth 2 --concurrency 5 --timeout 15
```

## Flags
`--url`: The URL to start crawling from.<br>
`--depth`: The maximum depth to crawl. Default is 3.<br>
`--concurrency`: The number of concurrent requests to make. Default is 5.<br>
`--timeout`: HTTP request timeout in seconds. Default is 10.

## Test
```bash
go test ./...
```
