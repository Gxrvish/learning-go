# Assignment 1: Log Line Parser CLI (`loggrep`)

**Tier:** 1 (Week 1-2) · **Time:** 6-8h · **Deps:** stdlib only

## 🎯 Learning Objectives
- Go syntax: variables, types, control flow, functions
- Error handling with multiple return values
- Standard library: `bufio.Scanner`, `os`, `flag`
- Idiomatic string parsing without regex
- Exit codes and CLI conventions

## 📋 Problem Statement
SRE team drowning in nginx access logs. Build CLI tool `loggrep` that reads access logs from stdin or file, filters by HTTP status code range, time window, and path prefix, then emits matching lines + summary stats to stderr. Runs in containerized log shippers — must handle 10GB files without OOM.

Sample input (Combined Log Format):
```
127.0.0.1 - alice [10/Oct/2023:13:55:36 +0000] "GET /api/users HTTP/1.1" 200 2326 "-" "curl/7.68.0"
```

## 📌 Requirements

**Functional:**
- Flags: `--status-min`, `--status-max`, `--path-prefix`, `--since`, `--until`, `--file` (default stdin)
- Parse Combined Log Format
- Matched lines → stdout unchanged
- Summary → stderr: total, matched, parse errors, p50/p95/p99 of response sizes
- Exit 0 if matches, 1 if none, 2 on fatal error

**Non-Functional:**
- ≥200 MB/s on modern laptop SSD
- O(1) memory regardless of file size
- No third-party dependencies

## ✅ Acceptance Criteria
- [ ] Streams 10GB file with <50MB RSS
- [ ] Handles malformed lines without crashing
- [ ] Time parsing handles Combined Log Format timezone offsets
- [ ] Correct exit codes
- [ ] `go vet` and `gofmt` clean

## 🔥 Advanced Challenge
Add `--follow` (like `tail -f`) using `os.File.Seek` + polling, with graceful Ctrl+C via `signal.Notify`.

## 🔍 Code Review Checklist
- [ ] No `panic` outside `main`
- [ ] Errors wrapped with `fmt.Errorf("...: %w", err)`
- [ ] `defer file.Close()` after error check on Open
- [ ] `bufio.Scanner` buffer sized for long lines (`scanner.Buffer`)
- [ ] No regex (force manual parsing)
- [ ] Receivers consistent (value vs pointer)
- [ ] `gofmt -s` clean

## 📊 Benchmarks
Target: `BenchmarkParseLine` <500 ns/op, ≤2 allocs/op.

## 💡 Production Insights
Inner loop of every log pipeline: Vector, Fluent Bit, Promtail. Allocation count matters: 1B lines × 1 alloc = 1B GC events.

## 🎓 Key Takeaways
1. Streaming > buffered for large inputs
2. `bufio.Scanner` default 64KB line limit — must raise
3. Errors are values, not exceptions
4. CLI exit codes carry semantics
5. Benchmark before optimizing

## 📚 Reference
- Effective Go: https://go.dev/doc/effective_go
- `pkg/bufio` source
- Nginx log format spec
