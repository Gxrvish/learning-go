package main

import "testing"

func TestParseLine(t *testing.T) {
	// TODO: table-driven test covering: happy path, missing fields,
	// quoted spaces in user-agent, malformed time, IPv6 remote addr.
	t.Skip("implement me")
}

func BenchmarkParseLine(b *testing.B) {
	line := `127.0.0.1 - alice [10/Oct/2023:13:55:36 +0000] "GET /api/users HTTP/1.1" 200 2326 "-" "curl/7.68.0"`
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = ParseLine(line)
	}
}
