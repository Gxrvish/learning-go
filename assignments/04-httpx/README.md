# Assignment 4: Resilient HTTP Client (`httpx`)

**Tier:** 4 (Week 5) · **Time:** 10-12h · **Deps:** stdlib, optional `golang.org/x/time/rate`

## 🎯 Learning Objectives
- Error wrapping (`errors.Is`, `errors.As`, `%w`)
- Custom + sentinel errors
- `log/slog` structured logs (Go 1.21+)
- Retry semantics, idempotency
- Functional options pattern

## 📋 Problem Statement
Backend hammers flaky upstream APIs. Build `httpx` — retrying HTTP client with circuit breaker, jittered exponential backoff, structured logs, sentinel errors. No retry on non-idempotent by default. Respect `Retry-After`.

## 📌 Requirements

**Functional:**
- `Client.Do(ctx, req)` retries 5xx, 429, network errors
- Options: `WithMaxRetries`, `WithBackoff`, `WithLogger`, `WithRetryPolicy`
- Sentinels: `ErrCircuitOpen`, `ErrMaxRetriesExceeded`, `ErrNonRetriable`
- Wrap underlying error (so `errors.Is(err, context.DeadlineExceeded)` works)
- Respect `Retry-After` (seconds or HTTP-date)

**Non-Functional:**
- Zero goroutine leaks
- Retries at `slog.LevelDebug`, outcome at `Info`/`Warn`

## ✅ Acceptance Criteria
- [ ] `errors.Is(err, ErrMaxRetriesExceeded)` after exhaustion
- [ ] `errors.As(err, &netErr)` retrieves underlying network error
- [ ] Default policy retries GET/HEAD/PUT/DELETE, not POST/PATCH
- [ ] Body re-readable across retries (`req.GetBody`)
- [ ] Tests use `httptest.Server` for failure simulation
- [ ] Deterministic jitter seed in tests

## 🔥 Advanced Challenge
Token-bucket via `x/time/rate`. Circuit breaker (closed/open/half-open) with thresholds. Prove no leak via `goleak`.

## 🔍 Code Review Checklist
- [ ] No silent error swallowing
- [ ] All errors wrapped
- [ ] Sentinels are package-level `var`
- [ ] Logger nil-safe (default `slog.Default()`)
- [ ] Body always closed (even on error)
- [ ] Context honored at every sleep (`ctx.Done()` not `time.After`)

## 📊 Benchmarks
- `BenchmarkDoNoRetry` <5 µs overhead vs `http.DefaultClient`
- Retry decision logic <100 ns

## 💡 Production Insights
Mirrors `hashicorp/go-retryablehttp` + `sony/gobreaker`. Every microservice fleet has its variant.

## 🎓 Key Takeaways
1. `%w` to wrap, `errors.Is`/`As` to unwrap
2. Sentinels = public branch API
3. Functional options scale
4. Idempotency drives retry policy, not status alone
5. `slog` is the new standard

## 📚 Reference
- https://go.dev/blog/go1.13-errors
- https://go.dev/blog/slog
- RFC 7231 §7.1.3
