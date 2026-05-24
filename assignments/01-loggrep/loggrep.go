// Package main implements the loggrep CLI.
//
// loggrep streams nginx Combined Log Format from a file or stdin, filters
// entries by status range, path prefix and time window, writes matches to
// stdout unchanged, and reports a summary (counts + p50/p95/p99 response
// sizes) to stderr.
//
// Exit codes:
//
//	0  at least one match
//	1  no matches
//	2  fatal error (cannot open file, invalid flag, etc.)
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

// clfTimeLayout is the Combined Log Format timestamp inside the [ ] block.
const clfTimeLayout = "02/Jan/2006:15:04:05 -0700"

// maxLineBytes raises bufio.Scanner past its 64 KiB default. Real-world
// access logs occasionally exceed this when user-agents or referers are
// pathological. 1 MiB is comfortably above anything observed in practice.
const maxLineBytes = 1 << 20

// LogEntry represents one parsed Combined Log Format line.
type LogEntry struct {
	RemoteAddr string
	User       string
	Time       time.Time
	Method     string
	Path       string
	Proto      string
	Status     int
	BytesSent  int64
	Referer    string
	UserAgent  string
}

// Filter is the set of predicates applied to each parsed entry.
// Zero-value bounds mean "no constraint" in that dimension.
type Filter struct {
	StatusMin, StatusMax int
	PathPrefix           string
	Since, Until         time.Time
}

// ParseLine parses a single Combined Log Format line.
// Returns (entry, true) on success, (zero, false) on parse error.
//
// Layout:
//
//	remote - user [time] "method path proto" status bytes "referer" "ua"
//
// Parsed by hand-rolled index walking; no regex, no Split allocations on
// the hot path. The returned LogEntry's string fields are sub-slices of
// the input — callers must copy if they need to outlive the input buffer
// (bufio.Scanner reuses its buffer between scans).
func ParseLine(line string) (LogEntry, bool) {
	var e LogEntry

	// remote addr — up to first space
	i := strings.IndexByte(line, ' ')
	if i < 0 {
		return LogEntry{}, false
	}
	e.RemoteAddr = line[:i]
	rest := line[i+1:]

	// ident — single token, ignored ("-")
	i = strings.IndexByte(rest, ' ')
	if i < 0 {
		return LogEntry{}, false
	}
	rest = rest[i+1:]

	// user — single token
	i = strings.IndexByte(rest, ' ')
	if i < 0 {
		return LogEntry{}, false
	}
	e.User = rest[:i]
	rest = rest[i+1:]

	// [time] — bracketed
	if len(rest) == 0 || rest[0] != '[' {
		return LogEntry{}, false
	}
	j := strings.IndexByte(rest, ']')
	if j < 0 {
		return LogEntry{}, false
	}
	t, err := time.Parse(clfTimeLayout, rest[1:j])
	if err != nil {
		return LogEntry{}, false
	}
	e.Time = t
	rest = rest[j+1:]
	if len(rest) < 2 || rest[0] != ' ' {
		return LogEntry{}, false
	}
	rest = rest[1:]

	// "method path proto" — quoted request line
	if rest[0] != '"' {
		return LogEntry{}, false
	}
	j = strings.IndexByte(rest[1:], '"')
	if j < 0 {
		return LogEntry{}, false
	}
	req := rest[1 : 1+j]
	// split request into 3 tokens by spaces
	s1 := strings.IndexByte(req, ' ')
	if s1 < 0 {
		return LogEntry{}, false
	}
	s2 := strings.IndexByte(req[s1+1:], ' ')
	if s2 < 0 {
		return LogEntry{}, false
	}
	e.Method = req[:s1]
	e.Path = req[s1+1 : s1+1+s2]
	e.Proto = req[s1+1+s2+1:]
	rest = rest[1+j+1:]
	if len(rest) < 2 || rest[0] != ' ' {
		return LogEntry{}, false
	}
	rest = rest[1:]

	// status — integer
	i = strings.IndexByte(rest, ' ')
	if i < 0 {
		return LogEntry{}, false
	}
	status, err := strconv.Atoi(rest[:i])
	if err != nil {
		return LogEntry{}, false
	}
	e.Status = status
	rest = rest[i+1:]

	// bytes sent — integer (or "-" for unknown)
	i = strings.IndexByte(rest, ' ')
	if i < 0 {
		return LogEntry{}, false
	}
	if rest[:i] == "-" {
		e.BytesSent = 0
	} else {
		b, err := strconv.ParseInt(rest[:i], 10, 64)
		if err != nil {
			return LogEntry{}, false
		}
		e.BytesSent = b
	}
	rest = rest[i+1:]

	// "referer"
	if len(rest) == 0 || rest[0] != '"' {
		return LogEntry{}, false
	}
	j = strings.IndexByte(rest[1:], '"')
	if j < 0 {
		return LogEntry{}, false
	}
	e.Referer = rest[1 : 1+j]
	rest = rest[1+j+1:]
	if len(rest) < 2 || rest[0] != ' ' {
		return LogEntry{}, false
	}
	rest = rest[1:]

	// "user-agent" — last field, may contain spaces but no unescaped quotes
	if rest[0] != '"' {
		return LogEntry{}, false
	}
	j = strings.IndexByte(rest[1:], '"')
	if j < 0 {
		return LogEntry{}, false
	}
	e.UserAgent = rest[1 : 1+j]

	return e, true
}

// Match reports whether the entry passes the filter.
// Zero-value bounds are treated as "no constraint".
func (f Filter) Match(e LogEntry) bool {
	if f.StatusMin != 0 && e.Status < f.StatusMin {
		return false
	}
	if f.StatusMax != 0 && e.Status > f.StatusMax {
		return false
	}
	if f.PathPrefix != "" && !strings.HasPrefix(e.Path, f.PathPrefix) {
		return false
	}
	if !f.Since.IsZero() && e.Time.Before(f.Since) {
		return false
	}
	if !f.Until.IsZero() && e.Time.After(f.Until) {
		return false
	}
	return true
}

// stats accumulates summary numbers for the run.
type stats struct {
	total       int64
	matched     int64
	parseErrors int64
	sizes       []int64 // response sizes of matched entries
}

// percentile returns the p-th percentile (0..100) of xs using
// nearest-rank. xs is sorted in place. Returns 0 if xs is empty.
func percentile(xs []int64, p float64) int64 {
	if len(xs) == 0 {
		return 0
	}
	if !sort.SliceIsSorted(xs, func(i, j int) bool { return xs[i] < xs[j] }) {
		sort.Slice(xs, func(i, j int) bool { return xs[i] < xs[j] })
	}
	rank := int((p/100.0)*float64(len(xs)-1) + 0.5)
	if rank < 0 {
		rank = 0
	}
	if rank >= len(xs) {
		rank = len(xs) - 1
	}
	return xs[rank]
}

// run is the testable entry point: takes streams + parsed config, returns
// the exit code. Keeping main thin and run() injectable makes integration
// testing straightforward.
func run(in io.Reader, out, errOut io.Writer, f Filter) int {
	scanner := bufio.NewScanner(in)
	scanner.Buffer(make([]byte, 64*1024), maxLineBytes)

	bw := bufio.NewWriter(out)
	defer bw.Flush()

	var s stats
	for scanner.Scan() {
		s.total++
		line := scanner.Text()
		e, ok := ParseLine(line)
		if !ok {
			s.parseErrors++
			continue
		}
		if !f.Match(e) {
			continue
		}
		s.matched++
		s.sizes = append(s.sizes, e.BytesSent)
		if _, err := bw.WriteString(line); err != nil {
			fmt.Fprintf(errOut, "write: %v\n", err)
			return 2
		}
		if err := bw.WriteByte('\n'); err != nil {
			fmt.Fprintf(errOut, "write: %v\n", err)
			return 2
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(errOut, "scan: %v\n", err)
		return 2
	}

	p50 := percentile(s.sizes, 50)
	p95 := percentile(s.sizes, 95)
	p99 := percentile(s.sizes, 99)
	fmt.Fprintf(errOut,
		"total=%d matched=%d parse_errors=%d p50=%d p95=%d p99=%d\n",
		s.total, s.matched, s.parseErrors, p50, p95, p99)

	if s.matched == 0 {
		return 1
	}
	return 0
}

// parseTimeFlag parses an optional RFC3339 flag value.
// Empty string yields a zero Time. Invalid values yield an error.
func parseTimeFlag(name, v string) (time.Time, error) {
	if v == "" {
		return time.Time{}, nil
	}
	t, err := time.Parse(time.RFC3339, v)
	if err != nil {
		return time.Time{}, fmt.Errorf("--%s: %w", name, err)
	}
	return t, nil
}

// openInput returns the reader for the configured input. Callers must
// close the returned io.Closer (which may be a no-op for stdin).
func openInput(path string) (io.ReadCloser, error) {
	if path == "" {
		return io.NopCloser(os.Stdin), nil
	}
	fh, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", path, err)
	}
	return fh, nil
}

func main() {
	statusMin := flag.Int("status-min", 0, "minimum status code (inclusive; 0 = unset)")
	statusMax := flag.Int("status-max", 0, "maximum status code (inclusive; 0 = unset)")
	pathPrefix := flag.String("path-prefix", "", "path prefix to match")
	since := flag.String("since", "", "earliest time to match (RFC3339)")
	until := flag.String("until", "", "latest time to match (RFC3339)")
	file := flag.String("file", "", "input file (default: stdin)")
	flag.Parse()

	sinceTime, err := parseTimeFlag("since", *since)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	untilTime, err := parseTimeFlag("until", *until)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	in, err := openInput(*file)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	defer in.Close()

	f := Filter{
		StatusMin:  *statusMin,
		StatusMax:  *statusMax,
		PathPrefix: *pathPrefix,
		Since:      sinceTime,
		Until:      untilTime,
	}

	code := run(in, os.Stdout, os.Stderr, f)
	os.Exit(code)
}
