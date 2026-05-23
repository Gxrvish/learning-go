// Package ring implements a fixed-capacity time-series ring buffer.
package ring

import "time"

// Sample is one timestamped value.
type Sample[T any] struct {
	Time  time.Time
	Value T
}

// Ring is a fixed-capacity ring buffer of Samples.
// Single-writer / multi-reader; concurrent writers require external sync.
type Ring[T any] struct {
	// TODO: buf []Sample[T], head, size, capacity
}

// New constructs a Ring with the given fixed capacity.
func New[T any](capacity int) *Ring[T] {
	// TODO
	return nil
}

// Push appends a sample, overwriting the oldest when full.
// Must be zero-allocation.
func (r *Ring[T]) Push(t time.Time, v T) {
	// TODO
}

// Len returns the current number of stored samples.
func (r *Ring[T]) Len() int {
	// TODO
	return 0
}

// Snapshot returns all stored samples, oldest first. One allocation.
func (r *Ring[T]) Snapshot() []Sample[T] {
	// TODO
	return nil
}

// Window returns samples with Time within the last d of the newest sample.
// TODO: use sort.Search; binary search over chronological order.
func (r *Ring[T]) Window(d time.Duration) []Sample[T] {
	// TODO
	return nil
}
