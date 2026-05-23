# Assignment 5: Parallel Web Crawler

**Tier:** 5 (Week 6-7) · **Time:** 12-16h · **Deps:** stdlib, `golang.org/x/net/html`

## 🎯 Learning Objectives
- Goroutines: spawn, lifecycle, leaks
- `sync.WaitGroup`, `sync.Mutex`, `sync.Map`
- Bounded concurrency via semaphore channel
- Race detector
- `go.uber.org/goleak`

## 📋 Problem Statement
Concurrent crawler. Seed URL + max depth + same-domain restriction. Fetch, parse links, recurse. Respect `robots.txt`. Bound concurrency. Produce sitemap + broken-link report.

## 📌 Requirements

**Functional:**
- BFS/DFS with depth limit
- Dedup URLs (visited set)
- Bounded concurrency (worker count configurable)
- Per-host rate limit
- JSON report: `{url, status, depth, parent, brokenLinks[]}`

**Non-Functional:**
- 10k-page site under 60s @ 50 workers
- No goroutine leak on cancel
- `go test -race` clean

## ✅ Acceptance Criteria
- [ ] Results streamed via channel
- [ ] Context cancel stops workers <1s
- [ ] `goleak` in `TestMain`
- [ ] Race-clean
- [ ] Visited set thread-safe
- [ ] Respects robots.txt `Disallow`

## 🔥 Advanced Challenge
Sharded map (16 shards) replacing mutex+map. Benchmark contention drop. `--resume` from on-disk state.

## 🔍 Code Review Checklist
- [ ] Every `go func()` documents output-closer
- [ ] No goroutine without lifecycle plan
- [ ] Channel ownership clear
- [ ] Semaphore released on error paths
- [ ] `WaitGroup.Add` before `go`
- [ ] Response bodies always `Body.Close()`

## 📊 Benchmarks
Speedup curve: 1, 4, 16, 64 workers. Document diminishing-return point.

## 💡 Production Insights
Pattern in `gospider`, `katana`, sitemap generators, scanners. Bounded concurrency + visited + politeness = universal triangle.

## 🎓 Key Takeaways
1. Channels for ownership, mutex for shared mutable state
2. `WaitGroup.Add` before goroutine
3. Cancellation is design, not afterthought
4. Race detector mandatory
5. Goroutine leaks compound silently

## 📚 Reference
- https://go.dev/blog/pipelines
- https://github.com/uber-go/goleak
