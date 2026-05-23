// Package orders implements the order placement service.
package orders

import (
	"context"
	"database/sql"
)

// Order is the aggregate persisted per placement.
type Order struct {
	ID             string
	IdempotencyKey string
	Items          []Item
	TotalCents     int64
}

// Item is one line item in an order.
type Item struct {
	SKU     string
	Qty     int
	UnitCts int64
}

// PlaceOrderRequest is the input for PlaceOrder.
type PlaceOrderRequest struct {
	IdempotencyKey string
	Items          []Item
}

// Service is the order placement service.
type Service struct {
	DB *sql.DB
	// TODO: repos
}

// PlaceOrder atomically reserves stock, persists the order, and writes the outbox event.
// Retries on serialization failure (PG SQLSTATE 40001).
//
// TODO: implement: tx, idempotency-key check, stock decrement, insert order+items,
// outbox row, commit, retry loop with bounded attempts.
func (s *Service) PlaceOrder(ctx context.Context, req PlaceOrderRequest) (*Order, error) {
	return nil, nil
}
