package offerwalls

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/solefaucet/btcwall-api/models"
	"github.com/twinj/uuid"
)

// PtcwallCallback handles kiwiwall callback
func (h OfferwallHandler) PtcwallCallback() gin.HandlerFunc {
	return func(c *gin.Context) {
		payload := struct {
			Amount float64 `form:"r" binding:"required"`
			Status int64   `form:"c" binding:"required,eq=1|eq=2"` // 1 success 2 chargeback
		}{}
		if err := c.BindWith(&payload, binding.Form); err != nil {
			logOfferwallCallback(models.OfferwallNamePtcwall, c, err)
			return
		}

		offer := offerFromContext(c)
		offer.OfferwallName = models.OfferwallNamePtcwall
		offer.TransactionID = uuid.NewV4().String()
		offer.Amount = int64(payload.Amount)

		if err := h.handleOfferCallback(offer, payload.Status == 2); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.String(http.StatusOK, "ok")
	}
}
