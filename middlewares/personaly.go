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

// PersonalyAuthRequired rejects request if client ip is not in the list and signature not match
func PersonalyAuthRequired(whitelistIPs []string, appHash, secretKey string) gin.HandlerFunc {
	ips := make(map[string]struct{})
	for _, v := range whitelistIPs {
		ips[v] = struct{}{}
	}

	return func(c *gin.Context) {
		if _, ok := ips[c.ClientIP()]; !ok {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		data := fmt.Sprintf("%v:%v:%v", c.Query("user_id"), appHash, secretKey)
		if sign := fmt.Sprintf("%x", md5.Sum([]byte(data))); sign != c.Query("signature") {
			httprequest, _ := httputil.DumpRequest(c.Request, true)
			logrus.WithFields(logrus.Fields{
				"event":          models.LogEventValidateOfferwallSignature,
				"offerwall":      "Personaly",
				"id_combination": c.Query("user_id"),
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
