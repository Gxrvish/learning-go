# Assignment 2: Time-Series Ring Buffer

**Tier:** 2 (Week 2-3) · **Time:** 8-10h · **Deps:** stdlib only

## 🎯 Learning Objectives
- Structs, methods, value vs pointer receivers
- Small interfaces (`fmt.Stringer`, `io.Writer`)
- Composition over inheritance
- Generics (Go 1.18+)
- Internal package layout

## 📋 Problem Statement
Reusable ring buffer for time-series samples. Monitoring agent holds last N seconds of metrics before scrape. Single-writer / multi-reader contract (document — no locks yet). Generic over sample type.

## 📌 Requirements

**Functional:**
- `Ring[T any]` with `Push(t time.Time, v T)`, `Snapshot() []Sample[T]`, `Window(d time.Duration) []Sample[T]`
- Snapshot chronological (oldest first)
- O(1) push, O(n) snapshot
- Implements `fmt.Stringer`

**Non-Functional:**
- Zero allocations in `Push`
- Snapshot allocates once (single slice)
- Fixed capacity at construction

## ✅ Acceptance Criteria
- [ ] Table-driven tests: empty / partial / full / wraparound
- [ ] `Push` 0 allocs/op
- [ ] `Window` uses binary search (`sort.Search`)
- [ ] Works for `int`, `float64`, custom structs

## 🔥 Advanced Challenge
`Aggregate(d, fn func([]T) T) []Sample[T]` downsampling, single allocation.

## 🔍 Code Review Checklist
- [ ] Receiver type consistent
- [ ] No exported fields where method suffices
- [ ] Doc comments start with identifier name
- [ ] Generic constraint minimal (`any` vs `comparable`)
- [ ] Tests use `t.Helper()` appropriately

## 📊 Benchmarks
- `BenchmarkPush`: <20 ns/op, 0 allocs
- `BenchmarkSnapshot`: O(n) confirmed

## 💡 Production Insights
Prometheus `client_golang` histograms, VictoriaMetrics tsdb. Bounded memory + recency.

## 🎓 Key Takeaways
1. Generics replace `interface{}` for containers
2. Fixed capacity > grow-on-append for predictable memory
3. Snapshot pattern decouples reader from writer lifetime
4. `sort.Search` is library-provided
5. Allocation profiling drives API design

## 📚 Reference
- https://go.dev/doc/tutorial/generics
- Prometheus tsdb internals
