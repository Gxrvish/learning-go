# Assignment 13: Distributed Key-Value Store (`tinykv`) — Capstone

**Tier:** 13 (Week 17-18) · **Time:** 40-60h · **Deps:** `hashicorp/raft` or `etcd-io/raft`, `grpc`, `pebble` or `bolt`

## 🎯 Learning Objectives
- Raft consensus integration
- State machine design
- Snapshots, log compaction
- Cluster membership changes
- Idempotency, linearizability, leader leases
- End-to-end production-grade Go

## 📋 Problem Statement
3-node distributed KV store with Raft. `Get`, `Put`, `Delete`, `CAS`, watch streams. Persistent log + snapshot. Survives node failures. Linearizable reads via leader read-index. Optional stale reads from followers.

## 📌 Requirements

**Functional:**
- gRPC API: `Get`, `Put`, `Delete`, `CompareAndSwap`, `Watch(stream)`
- Raft-replicated log applied to state machine
- Snapshot at log threshold; restore on startup
- Add/remove cluster member dynamically
- Linearizable read (leader read-index) + serializable opt-in
- Graceful shutdown, signal handling

**Non-Functional:**
- 3-node sustained 5k ops/s mixed
- Tolerate 1-node failure with no client errors (after election)
- Full-cluster restart recovers correctly
- No data loss with `fsync` on commit
- Observability: leader, term, applied index, log size

## 📁 Layout
```
tinykv/
├── cmd/tinykv/main.go
├── internal/
│   ├── server/      # gRPC handlers
│   ├── fsm/         # state machine
│   ├── store/       # pebble/bolt wrapper
│   ├── transport/   # raft transport
│   └── cluster/     # membership
├── api/v1/          # proto
├── deploy/          # docker-compose 3-node
└── test/            # integration, chaos
```

## ✅ Acceptance Criteria
- [ ] Jepsen-lite test: kill leader during writes, verify linearizability
- [ ] Snapshot + log compaction trigger correctly
- [ ] Membership change consistent
- [ ] Watches deliver in commit order, no gaps after reconnect
- [ ] Metrics: term, leader, applied index, snapshot size
- [ ] Graceful shutdown transfers leadership

## 🔥 Advanced Challenge
Range queries (sorted KV). Leases for TTL keys. Multi-Raft sharding. Replace library Raft with own Paxos.

## 🔍 Code Review Checklist
- [ ] FSM deterministic (no `time.Now`, no unseeded rand)
- [ ] All writes through Raft, never direct
- [ ] Snapshot serialization versioned
- [ ] Log compaction safe vs in-flight followers
- [ ] Read-index waits for committed index
- [ ] No goroutine leaks on shutdown
- [ ] gRPC handlers context-aware
- [ ] Observability comprehensive

## 📊 Benchmarks
Write p99 <20ms (3-node, local, fsync). Read p99 <5ms linearizable. Stale <1ms. Throughput vs cluster size.

## 💡 Production Insights
etcd in miniature. Same patterns: Consul, TiKV, CockroachDB internals.

## 🎓 Key Takeaways
1. Consensus is hard but boxable
2. FSM determinism non-negotiable
3. Linearizability ≠ availability (CAP)
4. Snapshots tame log growth, complicate replication
5. Observability separates "works" from "operable"

## 📚 Reference
- Raft paper: https://raft.github.io/raft.pdf
- etcd: github.com/etcd-io/etcd
- hashicorp/raft: github.com/hashicorp/raft
- Jepsen: https://jepsen.io
