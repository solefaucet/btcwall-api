package offerwalls

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/solefaucet/btcwall-api/models"
)

// PointClickTrackCallback handles pointclicktrack callback
func (h OfferwallHandler) PointClickTrackCallback() gin.HandlerFunc {
	return func(c *gin.Context) {
		payload := struct {
			Amount    float64 `form:"wallad_currency_amount" binding:"required"`
			OfferID   string  `form:"campaign_id" binding:"required"`
			OfferName string  `form:"campaign_name"`
			Status    string  `form:"status" binding:"required,eq=credited|eq=reversed"`
		}{}
		if err := c.BindWith(&payload, binding.Form); err != nil {
			logOfferwallCallback(models.OfferwallNamePointClickTrack, c, err)
			return
		}

		offer := offerFromContext(c)
		offer.OfferName = payload.OfferName
		offer.OfferwallName = models.OfferwallNamePointClickTrack
		offer.TransactionID = fmt.Sprintf("%v|%v", payload.OfferID, offer.UserID)
		offer.Amount = int64(payload.Amount)

		if err := h.handleOfferCallback(offer, payload.Status == "reversed"); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.String(http.StatusOK, "1")
	}
}
