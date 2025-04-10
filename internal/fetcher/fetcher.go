package fetcher

// PageFetcher defines an interface for fetching URLs and returning a list of strings or an error.
type PageFetcher interface {
    Fetch(url string) ([]string, error)
}
