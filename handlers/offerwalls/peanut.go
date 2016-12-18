package offerwalls

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/solefaucet/btcwall-api/models"
)

// PeanutCallback handles personaly callback
func (h OfferwallHandler) PeanutCallback() gin.HandlerFunc {
	return func(c *gin.Context) {
		payload := struct {
			Amount        float64 `form:"currencyAmt" binding:"required"`
			OfferName     string  `form:"offerTitle"`
			TransactionID string  `form:"transactionId"`
		}{}
		if err := c.BindWith(&payload, binding.Form); err != nil {
			logOfferwallCallback(models.OfferwallNamePeanut, c, err)
			return
		}

		offer := offerFromContext(c)
		offer.OfferName = payload.OfferName
		offer.OfferwallName = models.OfferwallNamePeanut
		offer.TransactionID = payload.TransactionID
		offer.Amount = int64(payload.Amount)

		if err := h.handleOfferCallback(offer, payload.Amount < 0); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.String(http.StatusOK, "1")
	}
}
