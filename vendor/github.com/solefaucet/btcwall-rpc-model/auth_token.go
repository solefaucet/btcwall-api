package rpcmodels

import "time"

// AuthToken model
type AuthToken struct {
	ID          int64     `db:"id" json:"id"`
	PublisherID int64     `db:"publisher_id" json:"publisher_id"`
	Email       string    `db:"email" json:"email"`
	AuthToken   string    `db:"auth_token" json:"auth_token"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}
