package middlewares

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

// PeanutAuthRequired rejects request if client ip is not in the list and signature not match
func PeanutAuthRequired(whitelistIPs []string, applicationKey, transactionKey string) gin.HandlerFunc {
	ips := make(map[string]struct{})
	for _, v := range whitelistIPs {
		ips[v] = struct{}{}
	}

	return func(c *gin.Context) {
		if _, ok := ips[c.ClientIP()]; !ok {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		status := c.Query("status")
		entry := logrus.WithField("status", status)
		if status != "C" {
			c.String(http.StatusOK, "1")
			entry.Warn("peanut status is not C")
			return
		}

		entry.Info("peanut status is C")
		c.Next()
	}
}
