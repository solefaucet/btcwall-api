package offerwalls

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/solefaucet/btcwall-api/models"
)

// AdscendCallback handles adscend callback
func (h OfferwallHandler) AdscendCallback() gin.HandlerFunc {
	return func(c *gin.Context) {
		payload := struct {
			Amount        float64 `form:"amount" binding:"required"`
			TransactionID string  `form:"tx_id" binding:"required"`
			OfferName     string  `form:"offer_name"`
		}{}
		if err := c.BindWith(&payload, binding.Form); err != nil {
			logOfferwallCallback(models.OfferwallNameAdscend, c, err)
			return
		}

		offer := offerFromContext(c)
		offer.OfferName = payload.OfferName
		offer.OfferwallName = models.OfferwallNameAdscend
		offer.TransactionID = payload.TransactionID
		offer.Amount = int64(payload.Amount)

		if err := h.handleOfferCallback(offer, payload.Amount <= 0); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.String(http.StatusOK, "1")
	}
}
