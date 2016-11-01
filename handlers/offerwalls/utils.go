package offerwalls

import (
	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/solefaucet/btcwall-api/models"
)

func logOfferwallCallback(offerwallName string, c *gin.Context, err error) {
	logrus.WithFields(logrus.Fields{
		"event":          models.LogEventParseOfferwallCallback,
		"offerwall_name": offerwallName,
		"query":          c.Request.URL.Query().Encode(),
		"error":          err.Error(),
	}).Error("fail to parse callback request")
}
