// Package main implements the loggrep CLI.
//
// TODO: wire up flags, open input (file or stdin), stream lines via bufio.Scanner,
// parse with ParseLine, filter via Filter.Match, write matches to stdout,
// summary stats to stderr, exit with documented codes.
package main

import (
	"time"
)

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
type Filter struct {
	StatusMin, StatusMax int
	PathPrefix           string
	Since, Until         time.Time
}

// ParseLine parses a single Combined Log Format line.
// Returns (entry, true) on success, (zero, false) on parse error.
//
// TODO: implement without regex. Index by spaces, brackets, quotes.
func ParseLine(line string) (LogEntry, bool) {
	return LogEntry{}, false
}

// Match reports whether the entry passes the filter.
//
// TODO: implement; treat zero-value bounds as "no constraint".
func (f Filter) Match(e LogEntry) bool {
	return false
}

func main() {
	// TODO
}
