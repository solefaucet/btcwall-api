package middlewares

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	rpcmodels "github.com/solefaucet/btcwall-rpc-model"
)

// PublisherAuthRequired _
func PublisherAuthRequired(getAuthToken func(string) (*rpcmodels.AuthToken, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		authToken, err := getAuthToken(c.Request.Header.Get("Auth-Token"))
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		if authToken == nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if authToken.CreatedAt.AddDate(0, 1, 0).Before(time.Now()) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("auth_token", authToken)
		c.Next()
	}
}
