package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// PointClickTrackAuthRequired rejects request if client ip is not in the list
func PointClickTrackAuthRequired(whitelistIPs []string) gin.HandlerFunc {
	ips := make(map[string]struct{})
	for _, v := range whitelistIPs {
		ips[v] = struct{}{}
	}

	return func(c *gin.Context) {
		if _, ok := ips[c.ClientIP()]; !ok {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Next()
	}
}
