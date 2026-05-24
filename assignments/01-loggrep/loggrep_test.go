package main

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

const goldenLine = `127.0.0.1 - alice [10/Oct/2023:13:55:36 +0000] "GET /api/users HTTP/1.1" 200 2326 "-" "curl/7.68.0"`

func TestParseLine(t *testing.T) {
	tests := []struct {
		name   string
		line   string
		wantOK bool
		check  func(t *testing.T, e LogEntry)
	}{
		{
			name:   "happy path",
			line:   goldenLine,
			wantOK: true,
			check: func(t *testing.T, e LogEntry) {
				if e.RemoteAddr != "127.0.0.1" {
					t.Errorf("RemoteAddr = %q", e.RemoteAddr)
				}
				if e.User != "alice" {
					t.Errorf("User = %q", e.User)
				}
				if e.Method != "GET" {
					t.Errorf("Method = %q", e.Method)
				}
				if e.Path != "/api/users" {
					t.Errorf("Path = %q", e.Path)
				}
				if e.Proto != "HTTP/1.1" {
					t.Errorf("Proto = %q", e.Proto)
				}
				if e.Status != 200 {
					t.Errorf("Status = %d", e.Status)
				}
				if e.BytesSent != 2326 {
					t.Errorf("BytesSent = %d", e.BytesSent)
				}
				if e.UserAgent != "curl/7.68.0" {
					t.Errorf("UserAgent = %q", e.UserAgent)
				}
				want := time.Date(2023, time.October, 10, 13, 55, 36, 0, time.UTC)
				if !e.Time.Equal(want) {
					t.Errorf("Time = %s, want %s", e.Time, want)
				}
			},
		},
		{
			name:   "ipv6 remote",
			line:   `::1 - - [10/Oct/2023:13:55:36 +0000] "GET / HTTP/1.1" 200 0 "-" "ua"`,
			wantOK: true,
			check: func(t *testing.T, e LogEntry) {
				if e.RemoteAddr != "::1" {
					t.Errorf("RemoteAddr = %q", e.RemoteAddr)
				}
			},
		},
		{
			name:   "dash bytes sent",
			line:   `127.0.0.1 - - [10/Oct/2023:13:55:36 +0000] "GET / HTTP/1.1" 304 - "-" "ua"`,
			wantOK: true,
			check: func(t *testing.T, e LogEntry) {
				if e.BytesSent != 0 {
					t.Errorf("BytesSent = %d, want 0", e.BytesSent)
				}
			},
		},
		{
			name:   "user-agent with spaces",
			line:   `127.0.0.1 - - [10/Oct/2023:13:55:36 +0000] "GET / HTTP/1.1" 200 10 "-" "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36"`,
			wantOK: true,
			check: func(t *testing.T, e LogEntry) {
				if !strings.HasPrefix(e.UserAgent, "Mozilla/5.0") {
					t.Errorf("UserAgent = %q", e.UserAgent)
				}
			},
		},
		{
			name:   "non-utc tz offset",
			line:   `127.0.0.1 - - [10/Oct/2023:13:55:36 +0530] "GET / HTTP/1.1" 200 10 "-" "ua"`,
			wantOK: true,
			check: func(t *testing.T, e LogEntry) {
				_, off := e.Time.Zone()
				if off != 5*3600+30*60 {
					t.Errorf("tz offset = %d", off)
				}
			},
		},
		{name: "empty", line: "", wantOK: false},
		{name: "missing time brackets", line: `127.0.0.1 - - 10/Oct/2023:13:55:36 +0000 "GET / HTTP/1.1" 200 10 "-" "ua"`, wantOK: false},
		{name: "malformed time", line: `127.0.0.1 - - [not-a-time] "GET / HTTP/1.1" 200 10 "-" "ua"`, wantOK: false},
		{name: "non-numeric status", line: `127.0.0.1 - - [10/Oct/2023:13:55:36 +0000] "GET / HTTP/1.1" XYZ 10 "-" "ua"`, wantOK: false},
		{name: "missing fields", line: `127.0.0.1 - -`, wantOK: false},
		{name: "request missing proto", line: `127.0.0.1 - - [10/Oct/2023:13:55:36 +0000] "GET /" 200 10 "-" "ua"`, wantOK: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			e, ok := ParseLine(tc.line)
			if ok != tc.wantOK {
				t.Fatalf("ok = %v, want %v", ok, tc.wantOK)
			}
			if ok && tc.check != nil {
				tc.check(t, e)
			}
		})
	}
}

func TestFilterMatch(t *testing.T) {
	base := LogEntry{
		Status:    200,
		Path:      "/api/users",
		Time:      time.Date(2023, 10, 10, 12, 0, 0, 0, time.UTC),
		BytesSent: 100,
	}
	tests := []struct {
		name string
		f    Filter
		e    LogEntry
		want bool
	}{
		{"zero filter passes", Filter{}, base, true},
		{"status in range", Filter{StatusMin: 200, StatusMax: 299}, base, true},
		{"status below", Filter{StatusMin: 300}, base, false},
		{"status above", Filter{StatusMax: 199}, base, false},
		{"path prefix match", Filter{PathPrefix: "/api/"}, base, true},
		{"path prefix miss", Filter{PathPrefix: "/static/"}, base, false},
		{"since before time", Filter{Since: base.Time.Add(-time.Hour)}, base, true},
		{"since after time", Filter{Since: base.Time.Add(time.Hour)}, base, false},
		{"until after time", Filter{Until: base.Time.Add(time.Hour)}, base, true},
		{"until before time", Filter{Until: base.Time.Add(-time.Hour)}, base, false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.f.Match(tc.e); got != tc.want {
				t.Errorf("Match = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestPercentile(t *testing.T) {
	if got := percentile(nil, 50); got != 0 {
		t.Errorf("empty p50 = %d", got)
	}
	xs := []int64{10, 20, 30, 40, 50, 60, 70, 80, 90, 100}
	if got := percentile(append([]int64(nil), xs...), 50); got != 50 && got != 60 {
		// nearest-rank: rank = round(0.5*9) = round(4.5) = 5 → xs[5] = 60
		t.Errorf("p50 = %d", got)
	}
	if got := percentile(append([]int64(nil), xs...), 100); got != 100 {
		t.Errorf("p100 = %d", got)
	}
	if got := percentile(append([]int64(nil), xs...), 0); got != 10 {
		t.Errorf("p0 = %d", got)
	}
}

func TestRunStreamingAndExitCodes(t *testing.T) {
	input := strings.Join([]string{
		`127.0.0.1 - - [10/Oct/2023:13:55:36 +0000] "GET /api/a HTTP/1.1" 200 100 "-" "ua"`,
		`127.0.0.1 - - [10/Oct/2023:13:55:37 +0000] "GET /api/b HTTP/1.1" 404 50 "-" "ua"`,
		`127.0.0.1 - - [10/Oct/2023:13:55:38 +0000] "GET /static/c HTTP/1.1" 200 200 "-" "ua"`,
		`malformed line here`,
		`127.0.0.1 - - [10/Oct/2023:13:55:39 +0000] "GET /api/d HTTP/1.1" 500 300 "-" "ua"`,
	}, "\n") + "\n"

	t.Run("matches present", func(t *testing.T) {
		var out, errOut bytes.Buffer
		f := Filter{PathPrefix: "/api/"}
		code := run(strings.NewReader(input), &out, &errOut, f)
		if code != 0 {
			t.Errorf("exit = %d, want 0", code)
		}
		matched := strings.Count(out.String(), "\n")
		if matched != 3 {
			t.Errorf("matched lines = %d, want 3", matched)
		}
		if !strings.Contains(errOut.String(), "total=5") {
			t.Errorf("stderr missing total=5: %q", errOut.String())
		}
		if !strings.Contains(errOut.String(), "parse_errors=1") {
			t.Errorf("stderr missing parse_errors=1: %q", errOut.String())
		}
		if !strings.Contains(errOut.String(), "matched=3") {
			t.Errorf("stderr missing matched=3: %q", errOut.String())
		}
	})

	t.Run("no matches", func(t *testing.T) {
		var out, errOut bytes.Buffer
		f := Filter{PathPrefix: "/nope/"}
		code := run(strings.NewReader(input), &out, &errOut, f)
		if code != 1 {
			t.Errorf("exit = %d, want 1", code)
		}
		if out.Len() != 0 {
			t.Errorf("stdout = %q, want empty", out.String())
		}
	})

	t.Run("status filter", func(t *testing.T) {
		var out, errOut bytes.Buffer
		f := Filter{StatusMin: 400, StatusMax: 499}
		code := run(strings.NewReader(input), &out, &errOut, f)
		if code != 0 {
			t.Errorf("exit = %d", code)
		}
		if strings.Count(out.String(), "\n") != 1 {
			t.Errorf("want 1 matched, got %q", out.String())
		}
	})

	t.Run("time window", func(t *testing.T) {
		var out, errOut bytes.Buffer
		f := Filter{
			Since: time.Date(2023, 10, 10, 13, 55, 38, 0, time.UTC),
			Until: time.Date(2023, 10, 10, 13, 55, 38, 0, time.UTC),
		}
		code := run(strings.NewReader(input), &out, &errOut, f)
		if code != 0 {
			t.Errorf("exit = %d", code)
		}
		if !strings.Contains(out.String(), "/static/c") {
			t.Errorf("want /static/c, got %q", out.String())
		}
	})
}

func BenchmarkParseLine(b *testing.B) {
	line := goldenLine
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = ParseLine(line)
	}
}
