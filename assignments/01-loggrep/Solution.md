# Solution: `loggrep`

A streaming CLI that filters nginx Combined Log Format access logs by status, path, and time, then reports counts plus p50/p95/p99 of response sizes. Stdlib only.

---

## 1. Architecture

```
            ┌────────────┐
 stdin ──►  │ openInput  │ ──► io.ReadCloser
 file ──►   └────────────┘
                 │
                 ▼
     ┌────────────────────────────┐
     │ bufio.Scanner (1 MiB buf)  │  one allocation, reused per line
     └────────────────────────────┘
                 │   string (sub-slice of scanner buffer)
                 ▼
     ┌────────────────────────────┐
     │ ParseLine (no regex)       │  index/slice walk; 0 alloc/op
     └────────────────────────────┘
                 │ LogEntry (string fields are sub-slices)
                 ▼
     ┌────────────────────────────┐
     │ Filter.Match               │  short-circuit predicates
     └────────────────────────────┘
                 │ pass
                 ├──► bufio.Writer (stdout)  ── original line
                 └──► stats { total, matched, parse_errors, sizes[] }
                                          │
                                          ▼ on EOF
                                  summary → stderr
                                  exit code → OS
```

Key separation of concerns:

| Symbol         | Job                                                   |
| -------------- | ----------------------------------------------------- |
| `ParseLine`    | Pure function: string → `(LogEntry, bool)`. No I/O.   |
| `Filter.Match` | Pure predicate over `LogEntry`. No I/O.               |
| `run`          | Wires reader/writer/filter. Returns exit code. Testable. |
| `main`         | Flag parsing, file opening, calls `run`, `os.Exit`.   |

`run` takes `io.Reader` / `io.Writer` instead of using `os.Stdin`/`os.Stdout` directly. This is why `TestRunStreamingAndExitCodes` can feed a `strings.Reader` and assert against `bytes.Buffer` — no temp files, no process fork.

---

## 2. What we learned

### 2.1 Streaming vs buffering

A 10 GB file does not fit in memory. `bufio.Scanner` reads at most one line into a reusable buffer, so memory is **O(longest line)**, not **O(file size)**. The reservoir for the percentile slice grows with **matched** entries — not with file size — which is the only deviation from strict O(1). For truly bounded memory, a t-digest or reservoir sample would replace the slice (see §5).

### 2.2 `bufio.Scanner` gotcha — 64 KiB default

Default `MaxScanTokenSize` is 64 KiB. Real-world access logs can blow past it (huge `User-Agent` strings, embedded JSON bodies in custom formats). We raised it explicitly:

```go
scanner.Buffer(make([]byte, 64*1024), maxLineBytes) // maxLineBytes = 1 MiB
```

First arg = initial buffer; second arg = max growth. Forget this and you get a silent `bufio.ErrTooLong` halting the scan.

### 2.3 Hand-rolled parser, zero allocations

The naive parser is `strings.Split(line, " ")` — concise, but allocates a `[]string` + N substring headers every line. At 1 B lines, that's 1 B garbage objects.

Instead we walk the line with `strings.IndexByte` and slice. **String slicing in Go does not allocate** — the resulting string shares the underlying bytes with the input. `strconv.Atoi` on a sub-slice also avoids allocation for typical lengths.

Benchmark confirms it:

```
BenchmarkParseLine-8   4,469,636   264.4 ns/op   0 B/op   0 allocs/op
```

Target was <500 ns and ≤2 allocs. We hit 264 ns and 0 allocs.

### 2.4 The hidden lifetime bug

The returned `LogEntry` holds string fields that are **sub-slices of the scanner's buffer**. The next `scanner.Scan()` overwrites that buffer. So if we ever stored `LogEntry`s across iterations, every field would silently rot.

In `run` we only:

1. Read fields immediately (`Filter.Match`)
2. Write the **original line** (also from the scanner buffer) before the next `Scan()`

This is safe by construction. If a future change wants to buffer entries (e.g. sort by time), it must copy strings explicitly — typically by calling `string([]byte(...))` or `strings.Clone`. The doc comment on `ParseLine` flags this.

### 2.5 Errors are values

No `panic`s outside fatal-init. Parse failures return `(zero, false)` — a normal control-flow signal. Open failures are wrapped:

```go
return nil, fmt.Errorf("open %s: %w", path, err)
```

`%w` preserves the chain so callers can `errors.Is` / `errors.As`. (We don't need that here, but it's the idiomatic default.)

### 2.6 CLI exit codes are an API

`grep` set the precedent: 0 = found, 1 = none, 2 = error. Shell pipelines depend on it (`loggrep ... && alert.sh`). We follow it exactly.

### 2.7 Time zones in CLF

CLF timestamps look like `[10/Oct/2023:13:55:36 +0530]`. Go's reference layout is:

```
02/Jan/2006:15:04:05 -0700
```

Once parsed, `time.Time` carries an offset; comparisons across zones work correctly via `Before`/`After` (which use the absolute instant, not wall-clock).

### 2.8 Percentiles: pick your algorithm

Three options were on the table:

| Method            | Memory | Exactness | Code |
| ----------------- | ------ | --------- | ---- |
| Sort-and-index    | O(n)   | Exact     | Trivial |
| Reservoir sample  | O(k)   | Approx (CI)| Easy |
| t-digest / HDR    | O(k)   | Tight     | Complex |

We picked sort-and-index. n = matched (not total), and matched ≪ total for selective filters. For a production tool emitting metrics every minute we'd swap in a t-digest. The point is to know **which knob you're turning** — exactness vs memory.

`nearest-rank` definition used: `rank = round((p/100) * (n-1))`. Matches the simplest of the [eight percentile methods Hyndman & Fan enumerate](https://en.wikipedia.org/wiki/Percentile#The_nearest-rank_method); good enough for SRE dashboards.

---

## 3. Internals — line by line tour of `ParseLine`

Input:

```
127.0.0.1 - alice [10/Oct/2023:13:55:36 +0000] "GET /api/users HTTP/1.1" 200 2326 "-" "curl/7.68.0"
```

| Step | Sentinel | Code idiom |
| ---- | -------- | ---------- |
| 1. remote | first space | `IndexByte(line, ' ')` |
| 2. ident | next space, drop | skip token |
| 3. user  | next space | `IndexByte` |
| 4. `[time]` | `[` then `]` | bracket pair → `time.Parse` |
| 5. `"req"` | `"` then `"` | inside: 3 tokens by space |
| 6. status | space | `strconv.Atoi` on sub-slice |
| 7. bytes  | space | special-case `"-"` → 0; else `ParseInt` |
| 8. `"ref"` | quote pair | sub-slice |
| 9. `"ua"` | quote pair | sub-slice (final field) |

Every step **shrinks `rest`** by re-slicing past the consumed segment. No `[]string`, no `strings.Fields`, no regex compile.

Failure mode at each step: returns `(LogEntry{}, false)`. Caller bumps `parse_errors` and moves on. **Malformed lines do not crash the program** — we tested this explicitly with `malformed line here` in the table-driven test.

---

## 4. Tests

```
TestParseLine              — 11 sub-tests: happy path, IPv6, dash-bytes,
                              UA with spaces, non-UTC tz, empty,
                              missing brackets, bad time, bad status,
                              short line, request missing proto
TestFilterMatch            — 10 sub-tests: zero filter, status bounds,
                              path prefix, since/until windows
TestPercentile             — empty, p0, p50, p100
TestRunStreamingAndExitCodes
                           — end-to-end: matches/no-matches/status/time
                              asserts exit code, stdout content,
                              stderr summary fields
BenchmarkParseLine         — 264 ns/op, 0 allocs/op
```

All four test functions pass; benchmark beats targets.

`run`'s `io.Reader` / `io.Writer` plumbing is what makes the e2e test cheap — no subprocess, no tempdir, no goroutine for stdin.

---

## 5. What's not done (deliberate)

- **`--follow` (tail -f)**: listed as advanced; would need `os.File.Seek(0, io.SeekEnd)`, polling loop, `signal.Notify` for graceful Ctrl+C. Different shape of program (signal-driven loop instead of EOF-driven).
- **Bounded-memory percentiles**: see §2.8. Sort-based exact percentiles use O(matched). For a real log shipper, replace `stats.sizes []int64` with a t-digest.
- **Concurrent parse**: at 264 ns/line one core does ~3.8 M lines/sec ≈ 200+ MB/s already. A worker pool would help only on slower CPUs or richer per-line work (e.g. geoip lookup). Premature here.

---

## 6. Concepts checklist (from the assignment)

- [x] Go syntax — vars, types, control flow, multi-return errors
- [x] Error handling — `(value, bool)` for parse, wrapped `error` for I/O
- [x] `bufio.Scanner` — including the buffer-size gotcha
- [x] `os` / `flag` — stdin vs file, RFC3339 time flags
- [x] Idiomatic string parsing — no regex, zero allocations
- [x] Exit codes — 0/1/2 per grep convention
- [x] `defer file.Close()` after error check
- [x] No `panic` outside fatal init
- [x] `go vet` + `gofmt -s` clean
- [x] Benchmark — 264 ns/op, 0 allocs/op (target <500 ns, ≤2 allocs)
- [x] Receivers consistent — `Filter` uses value receiver (small, immutable)
