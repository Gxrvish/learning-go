# Assignment 10: Streaming Data Pipeline (NDJSON ETL)

**Tier:** 10 (Week 13-14) · **Time:** 20-24h · **Deps:** stdlib, `golang.org/x/sync/errgroup`

## 🎯 Learning Objectives
- Pipeline pattern (stages connected by channels)
- Fan-out / fan-in
- `errgroup` for coordinated cancellation
- `sync/atomic`
- Memory-pooled buffers (`sync.Pool`)

## 📋 Problem Statement
NDJSON streaming ETL: read newline-delimited JSON from N files concurrently, validate, enrich (lookup), aggregate by key (windowed), write to sink. Backpressure, partial-failure handling, graceful shutdown. Replaceable stages.

## 📌 Requirements

**Functional:**
- Stages: Source → Decoder → Validator → Enricher → Windower → Sink
- Each stage parallelizable (configurable worker count)
- Any stage failure cancels pipeline (errgroup)
- Tumbling window by event time
- Per-stage metrics: in/out/error

**Non-Functional:**
- 1M events/s on 8-core (synthetic)
- Constant memory regardless of input
- Zero data loss on graceful shutdown
- Deadlock-free across all worker configs

## ✅ Acceptance Criteria
- [ ] Cancellation propagates <1s
- [ ] No goroutine leak (goleak)
- [ ] Counters use `atomic.Int64`
- [ ] `sync.Pool` shows allocation reduction
- [ ] Windowing correct under out-of-order events (event time)
- [ ] Chaos test: worker panic → recovered OR clean failure per policy

## 🔥 Advanced Challenge
Watermark-based windowing. Exactly-once sink with checkpoints. Throughput curves: 1 worker/stage vs many.

## 🔍 Code Review Checklist
- [ ] Each stage closes its output channel on exit
- [ ] `errgroup.Wait` called, error returned
- [ ] No shared mutable state between workers (channels only; atomics exempted)
- [ ] `sync.Pool.Put` returns zeroed buffer (`buf = buf[:0]`)
- [ ] Clock injected for tests (no raw `time.Now` in logic)

## 📊 Benchmarks
Throughput vs worker count per stage. Allocs/event with vs without pool. End-to-end p50/p99.

## 💡 Production Insights
Beam, Flink, Benthos, Vector — pipeline DAGs. Same shape in Go.

## 🎓 Key Takeaways
1. Pipeline = stages + channels + lifecycle
2. errgroup = ctx cancel + first-error wins
3. Backpressure natural with unbuffered channels
4. `sync.Pool` amortizes, doesn't eliminate
5. Atomic counters > mutex for counters

## 📚 Reference
- https://go.dev/blog/pipelines
- https://pkg.go.dev/golang.org/x/sync/errgroup
