package offerwalls

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/solefaucet/btcwall-api/models"
	rpcmodels "github.com/solefaucet/btcwall-rpc-model"
)

// AdscendCallback handles adscend callback
func (o OfferwallHandler) AdscendCallback() gin.HandlerFunc {
	return func(c *gin.Context) {
		payload := struct {
			Amount        float64 `form:"amount" binding:"required"`
			TransactionID string  `form:"tx_id" binding:"required"`
			OfferName     string  `form:"offer_name"` // NOTE: to fix
		}{}
		if err := c.BindWith(&payload, binding.Form); err != nil {
			logOfferwallCallback(models.OfferwallNameAdscend, c, err)
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
			OfferwallName: models.OfferwallNameAdscend,
			TransactionID: payload.TransactionID,
			Amount:        int64(payload.Amount),
		}

		if err := o.handleOfferCallback(offer, payload.Amount < 0); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.String(http.StatusOK, "1")
	}
}
