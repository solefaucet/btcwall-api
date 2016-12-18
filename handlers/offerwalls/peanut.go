package offerwalls

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/solefaucet/btcwall-api/models"
)

// PeanutCallback handles personaly callback
func (h OfferwallHandler) PeanutCallback() gin.HandlerFunc {
	return func(c *gin.Context) {
		payload := struct {
			Amount        string `form:"currencyAmt" binding:"required"`
			OfferName     string `form:"offerTitle"`
			TransactionID string `form:"transactionId"`
		}{}
		if err := c.BindWith(&payload, binding.Form); err != nil {
			logOfferwallCallback(models.OfferwallNamePeanut, c, err)
			return
		}
		amount, err := strconv.ParseFloat(strings.Replace(payload.Amount, ",", "", -1), 64)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			logOfferwallCallback(models.OfferwallNamePeanut, c, err)
			return
		}

		offer := offerFromContext(c)
		offer.OfferName = payload.OfferName
		offer.OfferwallName = models.OfferwallNamePeanut
		offer.TransactionID = payload.TransactionID
		offer.Amount = int64(amount)

		if err := h.handleOfferCallback(offer, amount < 0); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.String(http.StatusOK, "1")
	}
}
