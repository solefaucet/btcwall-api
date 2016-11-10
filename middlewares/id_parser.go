package middlewares

import (
	"encoding/base64"
	"encoding/json"

	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/solefaucet/btcwall-api/models"
)

type idPayload struct {
	PublisherID int64  `json:"publisher_id"`
	SiteID      int64  `json:"site_id"`
	UserID      int64  `json:"user_id"`
	TrackID     string `json:"track_id"`
}

// IDParserMiddleware parse publisher_id, site_id, user_id from query
func IDParserMiddleware(key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		data, err := base64.StdEncoding.DecodeString(c.Query(key))
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"event": models.LogEventParseIDCombination,
				"error": err.Error(),
			}).Info("cannot base64 decode id combination")
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		var payload idPayload
		if err := json.Unmarshal(data, &payload); err != nil {
			logrus.WithFields(logrus.Fields{
				"event":          models.LogEventParseIDCombination,
				"id_combination": string(data),
				"error":          err.Error(),
			}).Info("cannot parse id combination")
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		c.Set("publisher_id", payload.PublisherID)
		c.Set("site_id", payload.SiteID)
		c.Set("user_id", payload.UserID)
		c.Set("track_id", payload.TrackID)

		c.Next()
	}
}
