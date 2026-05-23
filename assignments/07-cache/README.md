# Assignment 7: Context-Aware Cache (TTL + Singleflight + LRU)

**Tier:** 7 (Week 9) · **Time:** 10-12h · **Deps:** stdlib, `golang.org/x/sync/singleflight`

## 🎯 Learning Objectives
- `context.Context` propagation rules
- Cancellation, deadlines, values (and their abuse)
- `singleflight` request coalescing
- TTL + LRU/LFU eviction
- Read-mostly: `sync.RWMutex` vs `atomic.Pointer`

## 📋 Problem Statement
Service fan-out to slow upstream (300ms p99). In-memory cache: TTL eviction, max-entries LRU, request coalescing (N concurrent misses for same key → 1 upstream call). Caller's deadline respected even while waiting on coalesced flight.

## 📌 Requirements

**Functional:**
- `Cache[K comparable, V any]` with `Get(ctx, key, loader func(ctx) (V, error)) (V, error)`
- TTL per entry (configurable, optional jitter)
- LRU at max size
- Background reaper for expired entries
- `Stats()` exposes hits, misses, evictions, in-flight, coalesced

**Non-Functional:**
- 95% hit path lock-free or RLock-only
- Reaper doesn't block reads
- Caller ctx cancel aborts wait but not flight

## ✅ Acceptance Criteria
- [ ] 1000 goroutines, same key, cold → loader called exactly once
- [ ] Caller ctx deadline returns `context.DeadlineExceeded` without canceling flight
- [ ] Stats accurate under contention
- [ ] No goroutine leak after Close
- [ ] LRU correct under concurrent access

## 🔥 Advanced Challenge
2-tier (`sync.Map` + LRU). Prometheus metric hooks. Stale-while-revalidate.

## 🔍 Code Review Checklist
- [ ] `context.Value` not used for cache state
- [ ] Loader gets context detached from caller
- [ ] Eviction doesn't hold lock during callbacks
- [ ] Reaper stops on `Close`
- [ ] Public methods take context where blocking

## 📊 Benchmarks
Hit path <50 ns/op. Miss-with-coalesce: 1 loader call regardless of N.

## 💡 Production Insights
groupcache, ristretto, bigcache. singleflight is THE thundering-herd answer. Google frontends, Cloudflare workers.

## 🎓 Key Takeaways
1. Context = request-scoped, not config
2. singleflight decouples N waiters from 1 producer
3. Background goroutines need explicit lifecycle
4. LRU = doubly-linked list + map; `container/list`
5. RWMutex shines when reads >> writes

## 📚 Reference
- https://pkg.go.dev/golang.org/x/sync/singleflight
- https://github.com/dgraph-io/ristretto
