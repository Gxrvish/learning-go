// Package crawler is a concurrent, bounded, polite web crawler.
package crawler

import (
	"context"
	"net/http"
)

// Crawler is a configured crawler instance.
type Crawler struct {
	MaxDepth    int
	Concurrency int
	SameDomain  bool
	HTTPClient  *http.Client
}

// Result is one crawl outcome streamed to the caller.
type Result struct {
	URL         string
	Status      int
	Depth       int
	Parent      string
	BrokenLinks []string
	Err         string
}

// Crawl returns a channel of Results discovered from seed.
// The channel is closed when the crawl completes or ctx is cancelled.
//
// TODO: implement with bounded worker semaphore, thread-safe visited set,
// per-host rate limit, robots.txt check, and graceful shutdown.
func (c *Crawler) Crawl(ctx context.Context, seed string) (<-chan Result, error) {
	return nil, nil
}
