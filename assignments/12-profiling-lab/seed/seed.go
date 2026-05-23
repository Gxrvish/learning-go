// Package seed contains intentionally-inefficient code for the profiling lab.
// Find the perf bugs via profiling, fix them, and prove improvement with benchstat.
//
// Seeded issues (do NOT remove until profiled):
//   1. SlowParse unmarshals the entire payload to read a single field.
//   2. BigStruct is passed by value through several layers.
//   3. SharedCache uses a global mutex around a read-mostly map.
//   4. HandleRequest spawns an unbounded goroutine per call.
package seed

import (
	"encoding/json"
	"sync"
)

// BigStruct is intentionally 2KB+ to surface copy costs.
type BigStruct struct {
	Buf [2048]byte
	Tag string
}

// SlowParse pulls only `name` out of the JSON but unmarshals everything.
// TODO: replace with a streaming or targeted decoder.
func SlowParse(payload []byte) (string, error) {
	var m map[string]any
	if err := json.Unmarshal(payload, &m); err != nil {
		return "", err
	}
	name, _ := m["name"].(string)
	return name, nil
}

// PassByValue ferries BigStruct through layers by value. TODO: pointer or arena.
func PassByValue(b BigStruct) BigStruct { return b }

// SharedCache: read-mostly map behind a write-lock.
type SharedCache struct {
	mu sync.Mutex
	m  map[string]string
}

// Get holds a write lock for a pure read. TODO: RWMutex or sync.Map.
func (c *SharedCache) Get(k string) string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.m[k]
}

// HandleRequest spawns a goroutine per call without bound. TODO: worker pool.
func HandleRequest(work func()) {
	go work()
}
