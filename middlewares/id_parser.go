package middlewares

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	multierror "github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/solefaucet/btcwall-api/models"
)

// IDParserMiddleware parse publisher_id, site_id, user_id from query
func IDParserMiddleware(key, offerwallName string, fs ...func(string) string) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Query(key)
		for _, f := range fs {
			id = f(id)
		}
		publisherID, siteID, userID, trackID, _ := parseIDCombination(id, offerwallName)

		c.Set("publisher_id", publisherID)
		c.Set("site_id", siteID)
		c.Set("user_id", userID)
		c.Set("track_id", trackID)

		c.Next()
	}
}

func parseIDCombination(idCombination, offerwallName string) (publisherID, siteID, userID int64, trackID string, err error) {
	entry := logrus.WithFields(logrus.Fields{
		"event":          models.LogEventParseIDCombination,
		"offerwall_name": offerwallName,
		"id_combination": idCombination,
	})

	fields := strings.Split(idCombination, "_")
	if len(fields) != 4 {
		err = fmt.Errorf("invalid id combination format %v", idCombination)
		entry.WithField("error", err.Error()).Warn("invalid id combination format")
		return
	}

	errs := &multierror.Error{}
	appendErrorIfNotNil := func(err error) {
		if err != nil {
			errs = multierror.Append(errs, err)
		}
	}

	publisherID, err = strconv.ParseInt(fields[0], 10, 64)
	appendErrorIfNotNil(errors.Wrapf(err, "cannot parse publisher_id %v", fields[0]))

	siteID, err = strconv.ParseInt(fields[1], 10, 64)
	appendErrorIfNotNil(errors.Wrapf(err, "cannot parse site_id %v", fields[1]))

	userID, err = strconv.ParseInt(fields[2], 10, 64)
	appendErrorIfNotNil(errors.Wrapf(err, "cannot parse user_id %v", fields[2]))

	trackID = fields[3]

	entry = entry.WithFields(logrus.Fields{
		"publisher_id": publisherID,
		"site_id":      siteID,
		"user_id":      userID,
		"track_id":     trackID,
	})

	err = errs.ErrorOrNil()
	if err != nil {
		entry.WithField("error", err.Error()).Warn("cannot parse id combination")
		return
	}

	entry.Info("succedd to parse id combination")
	return
}
