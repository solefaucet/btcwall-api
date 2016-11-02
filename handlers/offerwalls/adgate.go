package offerwalls

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/solefaucet/btcwall-api/models"
	rpcmodels "github.com/solefaucet/btcwall-rpc-model"
)

// AdgateCallback handles adgate callback
func (o OfferwallHandler) AdgateCallback() gin.HandlerFunc {
	return func(c *gin.Context) {
		payload := struct {
			Amount        float64 `form:"point_value" binding:"required"`
			TransactionID string  `form:"tx_id" binding:"required"`
			OfferName     string  `form:"offer_name"`
			Status        int64   `form:"status" binding:"required,eq=1|eq=0"` // 1 success 0 chargeback
		}{}
		if err := c.BindWith(&payload, binding.Form); err != nil {
			logOfferwallCallback(models.OfferwallNameAdgate, c, err)
			return
		}

		publisherID := c.MustGet("publisher_id").(int64)
		siteID := c.MustGet("site_id").(int64)
		userID := c.MustGet("user_id").(int64)

		offer := rpcmodels.Offer{
			PublisherID:   publisherID,
			SiteID:        siteID,
			UserID:        userID,
			OfferName:     payload.OfferName,
			OfferwallName: models.OfferwallNameAdgate,
			TransactionID: payload.TransactionID,
			Amount:        int64(payload.Amount),
		}

		if err := o.handleOfferCallback(offer, payload.Status == 0); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.String(http.StatusOK, "1")
	}
}
