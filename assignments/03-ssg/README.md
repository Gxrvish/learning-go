# Assignment 3: Static Site Generator (`ssg`)

**Tier:** 3 (Week 4) · **Time:** 12-16h · **Deps:** stdlib, optional `golang.org/x/sync/errgroup`, `fsnotify`

## 🎯 Learning Objectives
- Module layout: `cmd/`, `internal/`, `pkg/`
- Package boundaries, exported API surface
- `io.Reader`/`io.Writer` composition
- Subtests, golden files, `testdata/`
- `go test -cover`

## 📋 Problem Statement
Build `ssg` — Hugo-lite. Reads Markdown from `content/`, applies `html/template` templates, writes static site to `public/`. Front matter (YAML), nested dirs preserved, asset copy.

## 📌 Requirements

**Functional:**
- Parse YAML front matter
- Markdown → HTML (write minimal parser yourself: headings, paragraphs, code blocks, links, emphasis — no MD lib for core)
- Apply `html/template` with page data
- Walk `content/` preserving structure
- `--watch` rebuild on file change (`fsnotify` allowed)

**Non-Functional:**
- 1000 pages build <2s
- Parallel render via goroutines
- Test coverage ≥85% on `internal/markdown`, `internal/site`

## 📁 Layout
```
ssg/
├── cmd/ssg/main.go
├── internal/
│   ├── markdown/
│   ├── frontmatter/
│   └── site/
├── pkg/render/
├── testdata/{input,golden}/
└── go.mod
```

## ✅ Acceptance Criteria
- [ ] `go test ./... -cover` ≥85%
- [ ] Golden file tests for markdown corpus
- [ ] `internal/` unimportable externally
- [ ] CLI: `ssg build`, `ssg watch`
- [ ] Idempotent: rebuild byte-identical

## 🔥 Advanced Challenge
`--draft` flag. `sitemap.xml`. Concurrent render with `errgroup`, prove linear speedup vs `GOMAXPROCS`.

## 🔍 Code Review Checklist
- [ ] No circular imports
- [ ] `internal/` used correctly
- [ ] Public API has `Example*` tests
- [ ] Golden update via `-update` flag pattern
- [ ] Template errors surface with file:line
- [ ] `html/template` not `text/template`

## 📊 Benchmarks
- `BenchmarkRenderPage` <100 µs/op
- Full-corpus build: ns/page

## 💡 Production Insights
Hugo, Jekyll, Eleventy share this shape. `internal/` discipline matches Kubernetes, etcd. Golden tests are Go standard for renderers.

## 🎓 Key Takeaways
1. `internal/` is enforced
2. Golden files > brittle string asserts for HTML
3. Interfaces small + consumer-side
4. `html/template` auto-escapes per context
5. Coverage is floor not goal

## 📚 Reference
- https://go.dev/doc/modules/layout
- https://go.dev/blog/examples
- Hugo source
