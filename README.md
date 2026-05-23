# learning-go

Progressive Go assignments: zero → production-ready.

## Structure

Each assignment lives under `assignments/NN-slug/` with:
- `README.md` — full spec (problem, requirements, acceptance criteria, benchmarks)
- Starter `.go` files — skeleton interfaces / types
- `go.mod` — initialize per-assignment (`cd assignments/NN-slug && go mod init <module>`)

## Tiers

| # | Tier | Slug | Focus |
|---|------|------|-------|
| 1 | 1 | [01-loggrep](assignments/01-loggrep) | Syntax, stdlib, CLI, streaming I/O |
| 2 | 2 | [02-ringbuffer](assignments/02-ringbuffer) | Structs, methods, generics |
| 3 | 3 | [03-ssg](assignments/03-ssg) | Packages, modules, testing, golden files |
| 4 | 4 | [04-httpx](assignments/04-httpx) | Errors, slog, retries, sentinels |
| 5 | 5 | [05-crawler](assignments/05-crawler) | Goroutines, sync, race detector |
| 6 | 6 | [06-workerpool](assignments/06-workerpool) | Channels, backpressure, fan-out/in |
| 7 | 7 | [07-cache](assignments/07-cache) | Context, singleflight, TTL/LRU |
| 8 | 8 | [08-tinyurl](assignments/08-tinyurl) | HTTP server, middleware, REST |
| 9 | 9 | [09-orders](assignments/09-orders) | Postgres, tx, migrations, idempotency |
| 10 | 10 | [10-pipeline](assignments/10-pipeline) | Pipelines, errgroup, sync.Pool |
| 11 | 11 | [11-grpc-chat](assignments/11-grpc-chat) | gRPC, streaming, interceptors |
| 12 | 12 | [12-profiling-lab](assignments/12-profiling-lab) | pprof, escape analysis, benchstat |
| 13 | 13 | [13-tinykv](assignments/13-tinykv) | Raft, distributed KV (capstone) |
| 14 | 13+ | [14-hardening](assignments/14-hardening) | Ops, metrics, tracing, k8s |

## Workflow

```bash
cd assignments/01-loggrep
go mod init github.com/garvish/learning-go/loggrep   # first time
go test ./... -race -cover
go test -bench=. -benchmem
```

## Rules

- No skipping benchmarks
- No skipping `-race`
- No skipping code-review checklist in each README
- Write tests before / alongside implementation
