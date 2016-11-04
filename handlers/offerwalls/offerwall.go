package offerwalls

import (
	"github.com/Sirupsen/logrus"
	"github.com/solefaucet/btcwall-api/models"
	rpcmodels "github.com/solefaucet/btcwall-rpc-model"
)

// OfferwallHandler _
type OfferwallHandler struct {
	offerwallWriter offerwallWriter
}

type offerwallWriter interface {
	CreateOffer(offer rpcmodels.Offer) (duplicated bool, err error)
	ChargebackOffer(offer rpcmodels.Offer) (alreadyChargeback bool, err error)
}

// New create a offerwall handler
func New(offerwallWriter offerwallWriter) OfferwallHandler {
	return OfferwallHandler{
		offerwallWriter: offerwallWriter,
	}
}

func (o OfferwallHandler) handleOfferCallback(offer rpcmodels.Offer, isChargeback bool) error {
	entry := logrus.WithFields(logrus.Fields{
		"event":          models.LogEventHandleOfferwallCallback,
		"publisher_id":   offer.PublisherID,
		"site_id":        offer.SiteID,
		"user_id":        offer.UserID,
		"is_chargeback":  isChargeback,
		"offerwall_name": offer.OfferwallName,
		"amount":         offer.Amount,
		"offer_name":     offer.OfferName,
		"transaction_id": offer.TransactionID,
	})

	switch isChargeback {
	case true:
		if _, err := o.offerwallWriter.ChargebackOffer(offer); err != nil {
			entry.WithField("error", err.Error()).Error("fail to chargeback offer")
			return err
		}

	case false:
		if _, err := o.offerwallWriter.CreateOffer(offer); err != nil {
			entry.WithField("error", err.Error()).Error("fail to add offer")
			return err
		}
	}

	entry.Info("succeed to handle offer callback")
	return nil
}
