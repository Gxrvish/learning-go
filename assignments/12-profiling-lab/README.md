# Assignment 12: Profiling & Optimization Lab

**Tier:** 12 (Week 16) · **Time:** 12-16h · **Deps:** stdlib `net/http/pprof`, `runtime/pprof`, `benchstat`

## 🎯 Learning Objectives
- CPU, heap, allocs, mutex, block, goroutine profiles
- Reading flame graphs
- Escape analysis (`go build -gcflags=-m`)
- `benchstat` for statistically valid comparisons
- Optimization patterns: preallocate, reduce escape, inline

## 📋 Problem Statement
Inherit "slow" JSON-processing service. Codebase intentionally seeded with perf bugs: needless allocs, lock contention, large struct copies, escape-to-heap, unbounded goroutine spawn. Identify via profiling, fix, prove improvement with `benchstat`.

## 📌 Requirements

**Functional:**
- Same external behavior pre/post
- Comprehensive bench suite present
- Each fix documented with before/after numbers
- `benchstat` output in PR

**Non-Functional:**
- ≥5x CPU improvement
- ≥10x allocation reduction in hot path
- p99 latency ≥3x lower

## 📁 Layout (seeded)
```
seed/
├── internal/parser/    # full unmarshal for 1 field
├── internal/cache/     # global mutex on read-mostly map
├── internal/handler/   # unbounded goroutine spawn
└── internal/util/      # 2KB struct passed by value
```

## ✅ Acceptance Criteria
- [ ] CPU flame graph PNG attached
- [ ] Heap profile before + after
- [ ] `benchstat old.txt new.txt`: statistically significant (p<0.05)
- [ ] Escape-analysis output annotated for hot funcs
- [ ] Existing tests still pass
- [ ] Each commit references its motivating profile

## 🔥 Advanced Challenge
Try `encoding/json/v2` or `easyjson`. Hand-rolled parser for top-1 hot path. NUMA pinning.

## 🔍 Code Review Checklist
- [ ] Commits ↔ profile evidence
- [ ] No micro-opts without data
- [ ] Readability not sacrificed without commensurate gain
- [ ] `sync.Pool` only where allocs dominate
- [ ] Generics monomorphization considered

## 📊 Benchmarks
Provide `before.txt`, `after.txt` from `go test -bench=. -count=10`. Report via `benchstat`.

## 💡 Production Insights
Every senior Go engineer reads pprof flame graphs in their sleep. Cloudflare, Google SRE, Datadog blogs model this. Profile → hypothesize → fix → measure.

## 🎓 Key Takeaways
1. Don't guess, measure
2. Allocs dominate small-string CPU cost
3. `-gcflags=-m` reveals escape decisions
4. `benchstat` significance > eyeballing
5. Inline budget is real (~80 nodes)

## 📚 Reference
- https://go.dev/blog/pprof
- https://github.com/golang/go/wiki/CompilerOptimizations
- Dave Cheney: "High Performance Go"
