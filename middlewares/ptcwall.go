package middlewares

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

// PtcwallAuthRequired rejects request if client ip is not in the list
func PtcwallAuthRequired(whitelistIPs []string, postbackPassword string) gin.HandlerFunc {
	ips := make(map[string]struct{})
	for _, v := range whitelistIPs {
		ips[v] = struct{}{}
	}

	return func(c *gin.Context) {
		password := c.Query("pwd")

		if _, ok := ips[c.ClientIP()]; !ok || password != postbackPassword || postbackPassword == "" {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		// 1 cash 2 points
		if typ := c.Query("t"); typ != "2" {
			c.AbortWithStatus(http.StatusBadRequest)
			logrus.WithFields(logrus.Fields{
				"event": "validate ptcwall rewarding type",
				"type":  typ,
				"query": c.Request.URL.Query().Encode(),
			}).Error("ptcwall is not rewarding points")
			return
		}

		c.Next()
	}
}
