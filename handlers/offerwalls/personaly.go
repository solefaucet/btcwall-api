package offerwalls

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/solefaucet/btcwall-api/models"
	rpcmodels "github.com/solefaucet/btcwall-rpc-model"
)

// PersonalyCallback handles personaly callback
func (o OfferwallHandler) PersonalyCallback() gin.HandlerFunc {
	return func(c *gin.Context) {
		payload := struct {
			Amount    float64 `form:"amount" binding:"required,gt=0"`
			OfferName string  `form:"offer_name"`
			OfferID   string  `form:"offer_id"`
		}{}
		if err := c.BindWith(&payload, binding.Form); err != nil {
			logOfferwallCallback(models.OfferwallNamePersonaly, c, err)
			return
		}

		publisherID := c.MustGet("publisher_id").(int64)
		siteID := c.MustGet("site_id").(int64)
		userID := c.MustGet("user_id").(int64)
		trackID := c.MustGet("track_id").(string)

		transactionID := fmt.Sprintf("%v|%v", payload.OfferID, userID)
		offer := rpcmodels.Offer{
			PublisherID:   publisherID,
			SiteID:        siteID,
			UserID:        userID,
			TrackID:       trackID,
			OfferName:     payload.OfferName,
			OfferwallName: models.OfferwallNamePersonaly,
			TransactionID: transactionID,
			Amount:        int64(payload.Amount),
		}

		if err := o.handleOfferCallback(offer, payload.Amount < 0); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.String(http.StatusOK, "1")
	}
}
