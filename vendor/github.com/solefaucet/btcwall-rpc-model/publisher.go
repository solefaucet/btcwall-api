package rpcmodels

import (
	"encoding/json"
	"time"
)

// Publisher model
type Publisher struct {
	ID        int64     `db:"id" json:"id"`
	Email     string    `db:"email" json:"email"`
	Password  string    `db:"password" json:"-"`
	Address   string    `db:"address" json:"address"`
	Balance   int64     `db:"balance" json:"balance"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// MarshalJSON implements json marshaler
func (p Publisher) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"id":         p.ID,
		"email":      p.Email,
		"address":    p.Address,
		"balance":    float64(p.Balance) / 1e8,
		"created_at": p.CreatedAt,
		"updated_at": p.UpdatedAt,
	})
}
