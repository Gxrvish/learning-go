// Package pipeline implements a stage-based NDJSON ETL with cancellation and backpressure.
package pipeline

import "context"

// Sink consumes the final aggregated records.
type Sink interface {
	Write(ctx context.Context, key string, value any) error
	Close() error
}

// Run executes the pipeline over the given source files writing to sink.
// Returns on completion, fatal error, or ctx cancellation.
//
// TODO: build stages via errgroup, fan-out workers per stage, atomic metrics,
// sync.Pool for decode buffers, watermark/tumbling-window aggregator.
func Run(ctx context.Context, sources []string, sink Sink) error {
	return nil
}
