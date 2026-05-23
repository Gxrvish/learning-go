# Assignment 9: Order Service with Postgres + Migrations

**Tier:** 9 (Week 12) · **Time:** 16-20h · **Deps:** `database/sql`, `github.com/jackc/pgx/v5`, `github.com/pressly/goose`

## 🎯 Learning Objectives
- `database/sql` connection pool tuning
- Prepared statements
- Transactions, isolation levels, retry on serialization failure
- Migrations as code
- Repository pattern done right

## 📋 Problem Statement
E-commerce order service. Create order (multi-line items), reserve inventory, charge payment (mock), persist atomically. Idempotency-key support. Outbox pattern. Handles concurrent stock decrements without overselling.

## 📌 Requirements

**Functional:**
- `POST /v1/orders` accepts `Idempotency-Key` header, same key → same result (24h)
- Stock decrement: `SELECT ... FOR UPDATE` or optimistic version column
- Order + items + outbox event in single tx
- Background outbox publisher → stdout (mock)
- Reject if any item out of stock (atomic)

**Non-Functional:**
- No oversells under load: 100 buyers / 10 stock → exactly 10 succeed
- Connection pool explicit: `MaxOpen`, `MaxIdle`, `MaxLifetime` with justification
- Migrations versioned and reversible

## ✅ Acceptance Criteria
- [ ] Concurrent test proves no oversell
- [ ] Idempotency key returns identical response on retry
- [ ] Migrations up/down both work
- [ ] All queries use prepared statements or named args
- [ ] No SQL injection (no `fmt.Sprintf` into queries)
- [ ] Transient errors retried on 40001 (serialization)

## 🔥 Advanced Challenge
Outbox publisher via `LISTEN/NOTIFY`. Saga for distributed payment. Bulk via `pgx.CopyFrom`.

## 🔍 Code Review Checklist
- [ ] `rows.Close()` deferred on every Query
- [ ] `rows.Err()` checked after loop
- [ ] `tx.Rollback()` deferred (idempotent after Commit)
- [ ] No N+1
- [ ] Pool sized vs db `max_connections`
- [ ] Context on every DB call (`*Context` variants)
- [ ] `sql.ErrNoRows` → domain error at boundary

## 📊 Benchmarks
p99 order <50ms. Throughput ≥500/s single instance.

## 💡 Production Insights
Outbox: Shopify, Uber, every event-driven shop. `SELECT FOR UPDATE` vs optimistic = canonical trade-off. pgx is de facto Postgres driver.

## 🎓 Key Takeaways
1. Transactions need explicit retry on 40001
2. Idempotency keys = client-supplied dedup
3. Outbox = atomic state + events
4. Migrations live in repo, runnable in CI
5. Pool tuning is a DB concern

## 📚 Reference
- https://www.postgresql.org/docs/current/transaction-iso.html
- https://github.com/jackc/pgx
