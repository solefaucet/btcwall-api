package offerwalls

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/solefaucet/btcwall-api/models"
	rpcmodels "github.com/solefaucet/btcwall-rpc-model"
	"github.com/twinj/uuid"
)

// PtcwallCallback handles kiwiwall callback
func (o OfferwallHandler) PtcwallCallback() gin.HandlerFunc {
	return func(c *gin.Context) {
		payload := struct {
			Amount float64 `form:"r" binding:"required"`
			Status int64   `form:"c" binding:"required,eq=1|eq=2"` // 1 success 2 chargeback
		}{}
		if err := c.BindWith(&payload, binding.Form); err != nil {
			logOfferwallCallback(models.OfferwallNamePtcwall, c, err)
			return
		}

		publisherID := c.MustGet("publisher_id").(int64)
		siteID := c.MustGet("site_id").(int64)
		userID := c.MustGet("user_id").(int64)
		trackID := c.MustGet("track_id").(string)

		offer := rpcmodels.Offer{
			PublisherID:   publisherID,
			SiteID:        siteID,
			UserID:        userID,
			TrackID:       trackID,
			OfferName:     "", // NOTE: no offername provided by ptcwall
			OfferwallName: models.OfferwallNamePtcwall,
			TransactionID: uuid.NewV4().String(), // NOTE: no transaction_id provided by ptcwall
			Amount:        int64(payload.Amount),
		}

		if err := o.handleOfferCallback(offer, payload.Status == 2); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.String(http.StatusOK, "ok")
	}
}
