package seed

import (
	"sync"
	"testing"
)

var samplePayload = []byte(`{"name":"alice","age":30,"tags":["a","b","c"],"meta":{"k":"v"}}`)

func BenchmarkSlowParse(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = SlowParse(samplePayload)
	}
}

func BenchmarkPassByValue(b *testing.B) {
	bs := BigStruct{Tag: "x"}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		bs = PassByValue(bs)
	}
	_ = bs
}

func BenchmarkSharedCache(b *testing.B) {
	c := &SharedCache{m: map[string]string{"k": "v"}}
	b.ReportAllocs()
	var wg sync.WaitGroup
	for w := 0; w < 8; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < b.N/8; i++ {
				_ = c.Get("k")
			}
		}()
	}
	wg.Wait()
}
