package offerwalls

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/solefaucet/btcwall-api/models"
)

// AdgateCallback handles adgate callback
func (h OfferwallHandler) AdgateCallback() gin.HandlerFunc {
	return func(c *gin.Context) {
		payload := struct {
			Amount    float64 `form:"point_value" binding:"required"`
			OfferID   string  `form:"offer_id" binding:"required"`
			OfferName string  `form:"offer_name"`
			Status    int64   `form:"status" binding:"required,eq=1|eq=0"` // 1 success 0 chargeback
		}{}
		if err := c.BindWith(&payload, binding.Form); err != nil {
			logOfferwallCallback(models.OfferwallNameAdgate, c, err)
			return
		}

		offer := offerFromContext(c)
		offer.OfferName = payload.OfferName
		offer.OfferwallName = models.OfferwallNameAdgate
		offer.TransactionID = fmt.Sprintf("%v|%v", payload.OfferID, offer.UserID)
		offer.Amount = int64(payload.Amount)

		if err := h.handleOfferCallback(offer, payload.Status == 0); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.String(http.StatusOK, "1")
	}
}
