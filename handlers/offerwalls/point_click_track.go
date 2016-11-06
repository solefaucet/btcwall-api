package offerwalls

import (
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/solefaucet/btcwall-api/models"
	rpcmodels "github.com/solefaucet/btcwall-rpc-model"
)

// PointClickTrackCallback handles pointclicktrack callback
func (o OfferwallHandler) PointClickTrackCallback() gin.HandlerFunc {
	return func(c *gin.Context) {
		payload := struct {
			Amount         float64 `form:"commission" binding:"required"`
			OfferID        string  `form:"campaign_id" binding:"required"`
			SID            string  `form:"sid1" binding:"required"`
			OfferName      string  `form:"campaign_name"`
			Status         string  `form:"status" binding:"required,eq=credited|eq=reversed"`
			ReversalReason string  `form:"reversal_reason"`
		}{}
		if err := c.BindWith(&payload, binding.Form); err != nil {
			logOfferwallCallback(models.OfferwallNamePointClickTrack, c, err)
			return
		}

		publisherID := c.MustGet("publisher_id").(int64)
		siteID := c.MustGet("site_id").(int64)
		userID := c.MustGet("user_id").(int64)

		transactionID := fmt.Sprintf("%v|%v", payload.OfferID, userID)
		offer := rpcmodels.Offer{
			PublisherID:   publisherID,
			SiteID:        siteID,
			UserID:        userID,
			OfferName:     payload.OfferName,
			OfferwallName: models.OfferwallNamePointClickTrack,
			TransactionID: transactionID,
			Amount:        int64(payload.Amount),
		}

		isChargeback := payload.Status == "reversed"
		if isChargeback {
			logrus.WithFields(logrus.Fields{
				"event":     "chargeback",
				"offerwall": "point_click_track",
				"reason":    payload.ReversalReason,
			}).Info("chargeback point click track")
		}

		if err := o.handleOfferCallback(offer, isChargeback); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.String(http.StatusOK, "1")
	}
}
