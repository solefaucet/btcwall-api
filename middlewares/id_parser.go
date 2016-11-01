package middlewares

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	multierror "github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/solefaucet/btcwall-api/models"
)

// IDParserMiddleware parse publisher_id, site_id, user_id from query
func IDParserMiddleware(key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		idCombination := c.Query(key)
		publisherID, siteID, userID, err := parseIDCombination(idCombination)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		c.Set("publisher_id", publisherID)
		c.Set("site_id", siteID)
		c.Set("user_id", userID)

		c.Next()
	}
}

func parseIDCombination(idCombination string) (publisherID, siteID, userID int64, err error) {
	fields := strings.Split(idCombination, "_")
	if len(fields) != 3 {
		err = fmt.Errorf("invalid id combination format %v", idCombination)
		return
	}

	errs := &multierror.Error{}
	appendErrorIfNotNil := func(err error) {
		if err != nil {
			errs = multierror.Append(errs, nil)
		}
	}

	publisherID, err = strconv.ParseInt(fields[0], 10, 64)
	appendErrorIfNotNil(errors.Wrapf(err, "cannot parse publisher_id %v", fields[0]))

	siteID, err = strconv.ParseInt(fields[1], 10, 64)
	appendErrorIfNotNil(errors.Wrapf(err, "cannot parse site_id %v", fields[1]))

	userID, err = strconv.ParseInt(fields[2], 10, 64)
	appendErrorIfNotNil(errors.Wrapf(err, "cannot parse user_id %v", fields[2]))

	err = errs.ErrorOrNil()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"event":          models.LogEventParseIDCombination,
			"id_combination": idCombination,
			"publisher_id":   publisherID,
			"site_id":        siteID,
			"user_id":        userID,
			"error":          err.Error(),
		}).Info("cannot parse id combination")
	}

	return
}
