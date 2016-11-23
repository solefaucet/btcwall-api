package offerwalls

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/solefaucet/btcwall-api/models"
	rpcmodels "github.com/solefaucet/btcwall-rpc-model"
)

// WannadsCallback handles adgate callback
func (h OfferwallHandler) WannadsCallback() gin.HandlerFunc {
	return func(c *gin.Context) {
		payload := struct {
			Amount        float64 `form:"reward" binding:"required"`
			TransactionID string  `form:"transId" binding:"required"`
			OfferName     string  `form:"compaign_name"`
			Status        int64   `form:"status" binding:"required,eq=1|eq=2"` // 1 success 2 chargeback
		}{}
		if err := c.BindWith(&payload, binding.Form); err != nil {
			logOfferwallCallback(models.OfferwallNameWannads, c, err)
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
			OfferName:     payload.OfferName,
			OfferwallName: models.OfferwallNameWannads,
			TransactionID: payload.TransactionID,
			Amount:        int64(payload.Amount),
		}

		if err := h.handleOfferCallback(offer, payload.Status == 2); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.String(http.StatusOK, "OK")
	}
}
