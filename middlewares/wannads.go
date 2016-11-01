package middlewares

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/solefaucet/btcwall-api/models"
)

// WannadsAuthRequired rejects request if client ip is not in the list and signature not match
func WannadsAuthRequired(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		data := fmt.Sprintf("%v%v%v%v", c.Query("subId"), c.Query("transId"), c.Query("reward"), secretKey)
		if sign := fmt.Sprintf("%x", md5.Sum([]byte(data))); sign != c.Query("signature") {
			httprequest, _ := httputil.DumpRequest(c.Request, true)
			logrus.WithFields(logrus.Fields{
				"event":          models.LogEventValidateOfferwallSignature,
				"offerwall":      "Wannads",
				"id_combination": c.Query("subId"),
				"signature":      sign,
				"q_signature":    c.Query("signature"),
				"request":        string(httprequest),
			}).Error("signature not matched")
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Next()
	}
}
