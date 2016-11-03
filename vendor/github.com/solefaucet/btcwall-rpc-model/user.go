package rpcmodels

import (
	"encoding/json"
	"time"
)

// User model
type User struct {
	ID             int64     `db:"id" json:"id"`
	Address        string    `db:"address" json:"address"`
	Balance        int64     `db:"balance" json:"balance"`
	PendingBalance int64     `db:"pending_balance" json:"pending_balance"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
}

// MarshalJSON implements json marshaler
func (u User) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"id":              u.ID,
		"address":         u.Address,
		"balance":         float64(u.Balance) / 1e8,
		"pending_balance": float64(u.PendingBalance) / 1e8,
		"created_at":      u.CreatedAt,
		"updated_at":      u.UpdatedAt,
	})
}
