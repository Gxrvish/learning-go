package ring

import (
	"testing"
	"time"
)

func TestRing(t *testing.T) {
	// TODO: table-driven; cover empty, partial, full, wraparound.
	t.Skip("implement me")
}

func BenchmarkPush(b *testing.B) {
	r := New[float64](1024)
	now := time.Now()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		r.Push(now, float64(i))
	}
}
