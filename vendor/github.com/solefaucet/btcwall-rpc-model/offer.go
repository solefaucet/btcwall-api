package rpcmodels

import (
	"encoding/json"
	"time"
)

// offer status
const (
	OfferStatusPending = iota
	OfferStatusCharged
	OfferStatusChargeback
)

var offerStatusMapping = map[int]string{
	OfferStatusPending:    "pending",
	OfferStatusCharged:    "charged",
	OfferStatusChargeback: "chargeback",
}

// Offer model
type Offer struct {
	ID            int64     `db:"id"`
	PublisherID   int64     `db:"publisher_id"`
	SiteID        int64     `db:"site_id"`
	UserID        int64     `db:"user_id"`
	OfferName     string    `db:"offer_name"`
	OfferwallName string    `db:"offerwall_name"`
	TransactionID string    `db:"transaction_id"`
	TrackID       string    `db:"track_id"`
	Amount        int64     `db:"amount"`
	Status        int       `db:"status"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

// MarshalJSON implements json.Marshaler
func (o Offer) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"id":             o.ID,
		"publisher_id":   o.PublisherID,
		"site_id":        o.SiteID,
		"user_id":        o.UserID,
		"offer_name":     o.OfferName,
		"offerwall_name": o.OfferwallName,
		"transaction_id": o.TransactionID,
		"amount":         float64(o.Amount) / 1e8,
		"status":         offerStatusMapping[o.Status],
		"created_at":     o.CreatedAt,
		"updated_at":     o.UpdatedAt,
	})
}

// IsPending returns true if offer is pending
func (o Offer) IsPending() bool { return o.Status == OfferStatusPending }

// IsCharged returns true if offer is charged
func (o Offer) IsCharged() bool { return o.Status == OfferStatusCharged }

// IsChargeback returns true if offer is chargeback
func (o Offer) IsChargeback() bool { return o.Status == OfferStatusChargeback }
