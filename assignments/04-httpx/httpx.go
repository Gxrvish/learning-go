// Package httpx provides a retrying HTTP client with circuit breaker, jittered
// backoff, structured logging, and sentinel errors.
package httpx

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"
)

// Sentinel errors callers can branch on with errors.Is.
var (
	ErrCircuitOpen        = errors.New("httpx: circuit open")
	ErrMaxRetriesExceeded = errors.New("httpx: max retries exceeded")
	ErrNonRetriable       = errors.New("httpx: non-retriable response")
)

// RetryPolicy decides whether to retry and how long to wait before the next attempt.
type RetryPolicy func(attempt int, resp *http.Response, err error) (retry bool, delay time.Duration)

// Option configures a Client.
type Option func(*Client)

// WithMaxRetries caps the number of retries (default 3).
func WithMaxRetries(n int) Option { return func(c *Client) { /* TODO */ _ = c } }

// WithLogger sets the structured logger (default slog.Default()).
func WithLogger(l *slog.Logger) Option { return func(c *Client) { /* TODO */ _ = c } }

// WithRetryPolicy overrides the default retry decision function.
func WithRetryPolicy(p RetryPolicy) Option { return func(c *Client) { /* TODO */ _ = c } }

// Client wraps an http.Client with retry, backoff, logging, and circuit breaking.
type Client struct {
	// TODO: http.Client, retry policy, logger, backoff, breaker state
}

// New constructs a Client with the given options.
func New(opts ...Option) *Client {
	c := &Client{}
	for _, o := range opts {
		o(c)
	}
	return c
}

// Do executes the request with retry semantics. Context cancellation aborts all waits.
func (c *Client) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	// TODO
	return nil, errors.New("not implemented")
}
