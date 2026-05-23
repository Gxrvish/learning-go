# Assignment 8: Production REST API — URL Shortener (`tinyurl`)

**Tier:** 8 (Week 10-11) · **Time:** 20-25h · **Deps:** `net/http` (Go 1.22+ ServeMux), `log/slog`

## 🎯 Learning Objectives
- `net/http` server fundamentals
- Routing + middleware composition
- Request validation, content negotiation
- Graceful shutdown
- Observability: request ID, structured logs, metrics

## 📋 Problem Statement
`tinyurl` service. POST long URL → short code. GET short → 302 + click counter. API-key auth. Per-key rate limit. Health + readiness probes. OpenAPI 3 spec.

## 📌 Requirements

**Functional:**
- `POST /v1/links` `{url, customAlias?}` → `{shortCode, shortUrl, expiresAt}`
- `GET /{code}` → 302 + async click increment
- `GET /v1/links/{code}/stats` (auth) → clicks, lastAccessed, createdAt
- `DELETE /v1/links/{code}` (auth)
- `GET /healthz`, `/readyz`
- API key in `X-API-Key`

**Non-Functional:**
- p99 redirect <5ms (in-mem)
- Graceful shutdown drains in-flight (max 30s)
- Every request logged with request ID (via context)
- 429 on rate limit
- RFC 7807 problem+json errors

## 📁 Layout
```
tinyurl/
├── cmd/server/main.go
├── internal/
│   ├── api/
│   ├── store/
│   ├── service/
│   ├── middleware/
│   └── config/
├── api/openapi.yaml
└── go.mod
```

## ✅ Acceptance Criteria
- [ ] Middleware: recover → requestID → log → auth → ratelimit → handler
- [ ] Graceful shutdown on SIGTERM, 30s drain
- [ ] Integration tests via `httptest.NewServer`
- [ ] Validates inputs (URL parseable, http/https, <2048 chars)
- [ ] `*http.Server` constructed with timeouts (no `ListenAndServe` direct)
- [ ] `ReadHeaderTimeout`, `WriteTimeout`, `IdleTimeout` set

## 🔥 Advanced Challenge
Postgres `Store` with transactions. Prometheus `/metrics`. OpenTelemetry tracing.

## 🔍 Code Review Checklist
- [ ] Handlers thin, service layer holds logic
- [ ] No globals; DI via constructor
- [ ] Context through every layer
- [ ] Errors mapped to HTTP codes in one place
- [ ] Recover middleware prevents panic-to-client
- [ ] Timeouts everywhere
- [ ] No `interface{}` in handler signatures

## 📊 Benchmarks
`wrk -t4 -c100 -d30s` ≥20k RPS redirect. Per-req alloc <2KB.

## 💡 Production Insights
Shape of every Go microservice. Handler/service/store layering, middleware chain, structured logs, RFC 7807 — industry standard.

## 🎓 Key Takeaways
1. `*http.Server` for timeout control
2. Middleware is `func(http.Handler) http.Handler`
3. Context top→bottom, errors bottom→top
4. Graceful shutdown via `srv.Shutdown(ctx)`
5. RFC 7807 standardizes error responses

## 📚 Reference
- https://go.dev/blog/routing-enhancements
- RFC 7807
