package rpcmodels

import (
	"encoding/json"
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
	ID            int64     `db:"id" json:"id"`
	UserID        int64     `db:"user_id" json:"user_id"`
	Address       string    `db:"address" json:"address"`
	Amount        int64     `db:"amount" json:"amount"`
	Status        string    `db:"status" json:"status"`
	TransactionID string    `db:"transaction_id" json:"tx_id"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
}

// MarshalJSON implements json.Marshaler
func (w UserWithdrawal) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"id":             w.ID,
		"user_id":        w.UserID,
		"address":        w.Address,
		"amount":         float64(w.Amount) / 1e8,
		"transaction_id": w.TransactionID,
		"status":         w.Status,
		"created_at":     w.CreatedAt,
		"updated_at":     w.UpdatedAt,
	})
}

// PublisherWithdrawal model
type PublisherWithdrawal struct {
	ID            int64     `db:"id" json:"id"`
	PublisherID   int64     `db:"publisher_id" json:"publisher_id"`
	Address       string    `db:"address" json:"address"`
	Amount        int64     `db:"amount" json:"amount"`
	Status        string    `db:"status" json:"status"`
	TransactionID string    `db:"transaction_id" json:"tx_id"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
}

// MarshalJSON implements json.Marshaler
func (w PublisherWithdrawal) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"id":             w.ID,
		"publisher_id":   w.PublisherID,
		"address":        w.Address,
		"amount":         float64(w.Amount) / 1e8,
		"transaction_id": w.TransactionID,
		"status":         w.Status,
		"created_at":     w.CreatedAt,
		"updated_at":     w.UpdatedAt,
	})
}
