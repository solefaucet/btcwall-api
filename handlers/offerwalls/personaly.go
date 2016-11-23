package offerwalls

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/solefaucet/btcwall-api/models"
)

// PersonalyCallback handles personaly callback
func (h OfferwallHandler) PersonalyCallback() gin.HandlerFunc {
	return func(c *gin.Context) {
		payload := struct {
			Amount    float64 `form:"amount" binding:"required"`
			OfferName string  `form:"offer_name"`
			OfferID   string  `form:"offer_id"`
		}{}
		if err := c.BindWith(&payload, binding.Form); err != nil {
			logOfferwallCallback(models.OfferwallNamePersonaly, c, err)
			return
		}

		offer := offerFromContext(c)
		offer.OfferName = payload.OfferName
		offer.OfferwallName = models.OfferwallNamePersonaly
		offer.TransactionID = fmt.Sprintf("%v|%v", payload.OfferID, offer.UserID)
		offer.Amount = int64(payload.Amount)

		if err := h.handleOfferCallback(offer, payload.Amount < 0); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.String(http.StatusOK, "1")
	}
}
