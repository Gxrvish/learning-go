// Package fsm is the deterministic state machine applied to the Raft log.
// IMPORTANT: must be deterministic. No time.Now, no unseeded rand,
// no map iteration that depends on iteration order, no goroutines.
package fsm

import "io"

// FSM applies committed Raft log entries to local state and supports snapshots.
type FSM struct {
	// TODO: store (pebble/bolt) handle, applied index, revision counter
}

// Apply executes a single committed log entry. Returns an interpreter-defined result.
func (f *FSM) Apply(entry []byte) any {
	// TODO
	return nil
}

// Snapshot returns a point-in-time snapshot of the FSM state.
func (f *FSM) Snapshot() (Snapshot, error) {
	// TODO
	return nil, nil
}

// Restore replaces FSM state from a snapshot reader.
func (f *FSM) Restore(r io.ReadCloser) error {
	defer r.Close()
	// TODO
	return nil
}

// Snapshot is a serializable point-in-time view of FSM state.
type Snapshot interface {
	Persist(w io.Writer) error
	Release()
}
