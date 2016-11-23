package offerwalls

import (
	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/solefaucet/btcwall-api/models"
	rpcmodels "github.com/solefaucet/btcwall-rpc-model"
)

func logOfferwallCallback(offerwallName string, c *gin.Context, err error) {
	logrus.WithFields(logrus.Fields{
		"event":          models.LogEventParseOfferwallCallback,
		"offerwall_name": offerwallName,
		"query":          c.Request.URL.Query().Encode(),
		"error":          err.Error(),
	}).Error("fail to parse callback request")
}

func offerFromContext(c *gin.Context) rpcmodels.Offer {
	return rpcmodels.Offer{
		PublisherID: c.MustGet("publisher_id").(int64),
		SiteID:      c.MustGet("site_id").(int64),
		UserID:      c.MustGet("user_id").(int64),
		TrackID:     c.MustGet("track_id").(string),
	}
}
