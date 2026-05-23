// Package pool implements a generic bounded worker pool with backpressure.
package pool

import (
	"context"
	"errors"
)

// Sentinel errors callers branch on.
var (
	ErrQueueFull = errors.New("pool: queue full")
	ErrShutdown  = errors.New("pool: shutdown")
)

// Result is the outcome of one async job.
type Result[R any] struct {
	Value R
	Err   error
}

// Pool is a fixed-size worker pool with a bounded queue.
type Pool[T, R any] struct {
	// TODO: workers, jobs chan, fn, sync.Once for shutdown, ctx
}

// New constructs a pool with the given worker count, queue capacity, and handler fn.
func New[T, R any](workers, queue int, fn func(context.Context, T) (R, error)) *Pool[T, R] {
	// TODO
	return nil
}

// Submit enqueues a job synchronously, returning its result.
func (p *Pool[T, R]) Submit(ctx context.Context, job T) (R, error) {
	var zero R
	return zero, errors.New("not implemented")
}

// SubmitAsync enqueues a job and returns a channel that will receive exactly one Result.
func (p *Pool[T, R]) SubmitAsync(job T) <-chan Result[R] {
	// TODO
	return nil
}

// Shutdown drains the queue, returning when all workers exit or ctx fires.
func (p *Pool[T, R]) Shutdown(ctx context.Context) error {
	// TODO
	return nil
}

// ShutdownNow cancels in-flight work and returns pending un-started jobs.
func (p *Pool[T, R]) ShutdownNow() []T {
	// TODO
	return nil
}
