package offerwalls

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/solefaucet/btcwall-api/models"
)

// WannadsCallback handles adgate callback
func (h OfferwallHandler) WannadsCallback() gin.HandlerFunc {
	return func(c *gin.Context) {
		payload := struct {
			Amount    float64 `form:"reward" binding:"required"`
			OfferName string  `form:"campaign_name"`
			OfferID   string  `form:"campaign_id"`
			Status    int64   `form:"status" binding:"required,eq=1|eq=2"` // 1 success 2 chargeback
		}{}
		if err := c.BindWith(&payload, binding.Form); err != nil {
			logOfferwallCallback(models.OfferwallNameWannads, c, err)
			return
		}

		offer := offerFromContext(c)
		offer.OfferName = payload.OfferName
		offer.OfferwallName = models.OfferwallNameWannads
		offer.TransactionID = fmt.Sprintf("%v|%v", payload.OfferID, offer.UserID)
		offer.Amount = int64(payload.Amount)

		if err := h.handleOfferCallback(offer, payload.Status == 2); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.String(http.StatusOK, "OK")
	}
}
