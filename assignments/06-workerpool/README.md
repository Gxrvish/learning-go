# Assignment 6: Bounded Worker Pool with Backpressure

**Tier:** 6 (Week 8) · **Time:** 10-12h · **Deps:** stdlib

## 🎯 Learning Objectives
- Channel patterns: fan-out, fan-in, generator
- `select` with multiple cases, default, timeout
- `sync.Once`, `sync.Pool`
- Channel closing: who, when
- Backpressure design

## 📋 Problem Statement
Generic `Pool[T, R]` accepting jobs (T), producing results (R) via `func(ctx, T) (R, error)`. Bounded queue. If queue full, `Submit` blocks or returns `ErrQueueFull` per mode. Graceful shutdown drains or aborts.

## 📌 Requirements

**Functional:**
- `New[T,R](workers, queue int, fn) *Pool[T,R]`
- `Submit(ctx, T) (R, error)` synchronous
- `SubmitAsync(T) <-chan Result[R]` async
- `Shutdown(ctx)` drains or hits deadline
- `ShutdownNow()` cancels in-flight, returns pending

**Non-Functional:**
- Zero allocs per dispatch (post-warmup) — use `sync.Pool`
- Race-free
- Documented invariants

## ✅ Acceptance Criteria
- [ ] Submit honors ctx cancel
- [ ] Shutdown idempotent (`sync.Once`)
- [ ] No goroutine outlives Shutdown
- [ ] Stress: 1M jobs / 100 workers / leak-free / race-free
- [ ] Benchmarks for both submit modes

## 🔥 Advanced Challenge
Priority queue (high/low) with starvation prevention. Per-job timeout. Work-stealing.

## 🔍 Code Review Checklist
- [ ] Each channel has identified closer
- [ ] No `select { default: }` spin loops
- [ ] `Shutdown` callable concurrently with `Submit`
- [ ] No panic on submit-after-shutdown (return error)
- [ ] Generic constraints minimal

## 📊 Benchmarks
Submit overhead <500 ns/job. Compare with `ants`, `tunny`.

## 💡 Production Insights
Kubernetes' workqueue, Caddy handlers, Temporal worker. Backpressure separates pool from goroutine factory.

## 🎓 Key Takeaways
1. Bounded queue = backpressure
2. Creator usually closes
3. `sync.Once` for shutdown idempotency
4. `select` = concurrency Swiss army knife
5. Generic pool decouples T and R

## 📚 Reference
- https://go.dev/blog/pipelines
- k8s.io/client-go/util/workqueue
