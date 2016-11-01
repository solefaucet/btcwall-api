package rpcmodels

import "time"

// Site model
type Site struct {
	ID          int64     `db:"id" json:"id"`
	PublisherID int64     `db:"publisher_id" json:"publisher_id"`
	SiteName    string    `db:"site_name" json:"site_name"`
	SiteURL     string    `db:"site_url" json:"site_url"`
	Status      string    `db:"status" json:"status"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}
