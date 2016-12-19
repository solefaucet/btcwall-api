package offerwalls

import (
	"fmt"
	"unicode/utf8"

	"github.com/Sirupsen/logrus"
	"github.com/solefaucet/btcwall-api/models"
	rpcmodels "github.com/solefaucet/btcwall-rpc-model"
)

// OfferwallHandler _
type OfferwallHandler struct {
	offerwallWriter            offerwallWriter
	runcpaRevenueShareNotifier runcpaRevenueShareNotifier
}

type offerwallWriter interface {
	CreateOffer(offer rpcmodels.Offer) (duplicated bool, err error)
	ChargebackOffer(offer rpcmodels.Offer) (alreadyChargeback bool, err error)
}

type runcpaRevenueShareNotifier interface {
	CallbackRevenueShare(trackID string, amount float64)
}

// New create a offerwall handler
func New(offerwallWriter offerwallWriter, runcpaRevenueShareNotifier runcpaRevenueShareNotifier) OfferwallHandler {
	return OfferwallHandler{
		offerwallWriter:            offerwallWriter,
		runcpaRevenueShareNotifier: runcpaRevenueShareNotifier,
	}
}

func (h OfferwallHandler) handleOfferCallback(offer rpcmodels.Offer, isChargeback bool) error {
	entry := logrus.WithFields(logrus.Fields{
		"event":          models.LogEventHandleOfferwallCallback,
		"publisher_id":   offer.PublisherID,
		"site_id":        offer.SiteID,
		"user_id":        offer.UserID,
		"track_id":       offer.TrackID,
		"is_chargeback":  fmt.Sprint(isChargeback),
		"offerwall_name": offer.OfferwallName,
		"amount":         offer.Amount,
		"offer_name":     offer.OfferName,
		"transaction_id": offer.TransactionID,
	})

	// workaround for kiwiwall offer_name containing invalid characters
	offer.OfferName = validOffername(offer.OfferName)

	switch isChargeback {
	case true:
		if _, err := h.offerwallWriter.ChargebackOffer(offer); err != nil {
			entry.WithField("error", err.Error()).Error("fail to chargeback offer")
			return err
		}

	case false:
		duplicated, err := h.offerwallWriter.CreateOffer(offer)
		if err != nil {
			entry.WithField("error", err.Error()).Error("fail to add offer")
			return err
		}

		if !duplicated {
			h.runcpaRevenueShareNotifier.CallbackRevenueShare(offer.TrackID, float64(offer.Amount)/1e8)
		}
	}

	entry.Info("succeed to handle offer callback")
	return nil
}

func validOffername(name string) string {
	if utf8.ValidString(name) {
		return name
	}
	return ""
}
