# Assignment 11: gRPC Chat Service (Streaming + Auth)

**Tier:** 11 (Week 15) · **Time:** 16-20h · **Deps:** `google.golang.org/grpc`, `google.golang.org/protobuf`, `buf`

## 🎯 Learning Objectives
- protobuf schema design
- gRPC: unary, server-stream, bidi-stream
- Interceptors (gRPC middleware)
- mTLS or token auth
- Backwards-compat schema evolution

## 📋 Problem Statement
Real-time chat. Bidi for message exchange. Server-stream for presence. Unary for room mgmt. mTLS service-to-service. JWT for clients. Schema must evolve without breaking v1.

## 📌 Requirements

**Functional:**
- Unary: `CreateRoom`, `JoinRoom`, `LeaveRoom`
- Bidi: `Send(stream Message) returns (stream Message)`
- Server-stream: `WatchPresence(RoomID) returns (stream PresenceEvent)`
- Auth: JWT in metadata `authorization: Bearer ...`
- Interceptors: auth, log, recover, metrics

**Non-Functional:**
- 10k concurrent streams/server
- Graceful shutdown drains streams (30s max)
- Backwards-compat: adding fields doesn't break clients
- Deadlines respected on every RPC

## 📁 Layout
```
chat/
├── api/chat/v1/chat.proto
├── cmd/server/main.go
├── internal/{server,room,auth}/
└── buf.yaml, buf.gen.yaml
```

## ✅ Acceptance Criteria
- [ ] `buf lint` clean, `buf breaking` in CI
- [ ] mTLS between services (test with `crypto/tls.Config`)
- [ ] Streaming handlers clean up on disconnect
- [ ] Interceptors composable, ordered
- [ ] Load test: 1k concurrent bidi streams, leak-free
- [ ] Versioned package path (`chat.v1`)

## 🔥 Advanced Challenge
Server fan-out (Send → broadcast to room). Deadline propagation. `grpc.health.v1` service.

## 🔍 Code Review Checklist
- [ ] No `panic` in stream handlers
- [ ] Streams check `ctx.Err()` between sends
- [ ] Per-stream goroutines cleaned on disconnect
- [ ] Bounded send channels per client (backpressure)
- [ ] `status.Error` with appropriate `codes.*`
- [ ] Metadata extracted type-safely

## 📊 Benchmarks
Unary p99 <2ms in-cluster. Streaming ≥10k msg/s/stream. Idle stream <10KB RAM.

## 💡 Production Insights
gRPC = lingua franca: Kubernetes, etcd, Envoy, CockroachDB. Streaming is the underused superpower. buf is the modern toolchain.

## 🎓 Key Takeaways
1. Streams = stateful goroutines on server
2. Interceptors layer like HTTP middleware
3. Metadata = headers; trailers = post-call
4. Schema evolution: add only
5. Health check service is part of protocol

## 📚 Reference
- https://grpc.io/docs/languages/go/
- https://buf.build/docs
