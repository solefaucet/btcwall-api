package rpcmodels

import (
	"encoding/json"
	"fmt"
	"time"
)

// withdrawal status
const (
	WithdrawalStatusPending    = "pending"
	WithdrawalStatusProcessing = "processing"
	WithdrawalStatusProcessed  = "processed"
)

// UserWithdrawal model
type UserWithdrawal struct {
	ID            int64     `db:"id"`
	UserID        int64     `db:"user_id"`
	Address       string    `db:"address"`
	Amount        int64     `db:"amount"`
	Status        string    `db:"status"`
	TransactionID string    `db:"transaction_id"`
	UpdatedAt     time.Time `db:"updated_at"`
	CreatedAt     time.Time `db:"created_at"`
}

// MarshalJSON implements json.Marshaler
func (w UserWithdrawal) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"id":             w.ID,
		"user_id":        w.UserID,
		"address":        w.Address,
		"amount":         float64(w.Amount) / 1e8,
		"transaction_id": fmt.Sprintf("https://blockchain.info/tx/%s", w.TransactionID),
		"status":         w.Status,
		"created_at":     w.CreatedAt,
		"updated_at":     w.UpdatedAt,
	})
}

// PublisherWithdrawal model
type PublisherWithdrawal struct {
	ID            int64     `db:"id"`
	PublisherID   int64     `db:"publisher_id"`
	Address       string    `db:"address"`
	Amount        int64     `db:"amount"`
	Status        string    `db:"status"`
	TransactionID string    `db:"transaction_id"`
	UpdatedAt     time.Time `db:"updated_at"`
	CreatedAt     time.Time `db:"created_at"`
}

// MarshalJSON implements json.Marshaler
func (w PublisherWithdrawal) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"id":             w.ID,
		"publisher_id":   w.PublisherID,
		"address":        w.Address,
		"amount":         float64(w.Amount) / 1e8,
		"transaction_id": fmt.Sprintf("https://blockchain.info/tx/%s", w.TransactionID),
		"status":         w.Status,
		"created_at":     w.CreatedAt,
		"updated_at":     w.UpdatedAt,
	})
}
