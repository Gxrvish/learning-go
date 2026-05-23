// Package store defines the persistence interface for tinyurl.
package store

import (
	"context"
	"errors"
	"time"
)

// ErrNotFound is returned when no record matches the lookup.
var ErrNotFound = errors.New("store: not found")

// Link is a stored short-URL record.
type Link struct {
	Code        string
	URL         string
	CreatedAt   time.Time
	ExpiresAt   time.Time
	Clicks      int64
	LastAccessed time.Time
}

// Store is the persistence contract. Implementations: in-memory, Postgres, Redis.
type Store interface {
	Save(ctx context.Context, l Link) error
	Get(ctx context.Context, code string) (Link, error)
	IncrementClicks(ctx context.Context, code string) error
	Delete(ctx context.Context, code string) error
}
