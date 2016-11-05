package middlewares

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/solefaucet/btcwall-api/models"
)

// ProxyAuthRequired blocks request if it's from a proxy
func ProxyAuthRequired(getScoreByIP func(string) (int64, error), threshold int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		entry := logrus.WithFields(logrus.Fields{
			"event": models.LogEventBlockProxy,
			"ip":    ip,
			"block": false,
		})

		score, err := getScoreByIP(ip)
		entry = entry.WithField("score", score)

		if err != nil {
			entry.WithField("error", err.Error()).Warning("fail to get score by ip")
			c.Next()
			return
		}

		if score >= threshold {
			entry.WithField("block", true).Info("block request from proxy")
			c.AbortWithStatus(http.StatusProxyAuthRequired)
			return
		}

		entry.Info("request is not from proxy")
		c.Next()
	}
}
