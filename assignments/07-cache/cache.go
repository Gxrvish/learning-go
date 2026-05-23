// Package cache provides a generic TTL+LRU+singleflight cache.
package cache

import (
	"context"
	"time"
)

// Options configures a Cache.
type Options struct {
	MaxEntries int
	TTL        time.Duration
	Jitter     time.Duration
}

// Stats are aggregate cache counters.
type Stats struct {
	Hits, Misses, Evictions, InFlight, Coalesced int64
}

// Cache is a generic in-memory TTL+LRU cache with singleflight coalescing.
type Cache[K comparable, V any] struct {
	// TODO: lru list, map, singleflight.Group, reaper stop chan, stats atomics
}

// New constructs a Cache.
func New[K comparable, V any](opts Options) *Cache[K, V] {
	// TODO
	return nil
}

// Get returns the cached value, calling loader on miss. Concurrent calls for
// the same key share a single loader invocation.
func (c *Cache[K, V]) Get(ctx context.Context, key K, loader func(context.Context) (V, error)) (V, error) {
	var zero V
	return zero, nil
}

// Stats returns a snapshot of counters.
func (c *Cache[K, V]) Stats() Stats {
	return Stats{}
}

// Close stops the background reaper and releases resources.
func (c *Cache[K, V]) Close() error {
	return nil
}
