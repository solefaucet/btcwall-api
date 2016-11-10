package runcpa

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	multierror "github.com/hashicorp/go-multierror"
	"github.com/parnurzeal/gorequest"
	"github.com/pkg/errors"
	"github.com/solefaucet/btcwall-api/models"
)

// Notifier _
type Notifier struct {
	baseRegistrationCallbackURL string
	baseRevenueShareCallbackURL string

	buffer chan string
}

// New _
func New(baseRegistrationCallbackURL, baseRevenueShareCallbackURL string) Notifier {
	return Notifier{
		baseRegistrationCallbackURL: baseRegistrationCallbackURL,
		baseRevenueShareCallbackURL: baseRevenueShareCallbackURL,
	}
}

// CallbackRegistration _
func (n Notifier) CallbackRegistration(trackID string) {
	if trackID == "" {
		return
	}

	url := fmt.Sprintf("%s/%s", n.baseRegistrationCallbackURL, trackID)
	resp, _, errs := gorequest.New().Get(url).End()
	validStatusCode := resp.StatusCode >= 200 && resp.StatusCode < 300
	if len(errs) != 0 || validStatusCode {
		err := errors.Wrap(&multierror.Error{Errors: errs}, "fail to callback runcpa registration")
		logrus.WithFields(logrus.Fields{
			"event":            models.LogEventCallbackRuncpaRegistration,
			"url":              url,
			"http_status_code": resp.StatusCode,
			"error":            err.Error(),
		}).Warn("fail to callback runcpa registration")
	}
}

// CallbackRevenueShare _
func (n Notifier) CallbackRevenueShare(trackID string, sum float64) {
	if trackID == "" || sum <= 0 {
		return
	}

	url := fmt.Sprintf("%s/%s/%.8f", n.baseRevenueShareCallbackURL, trackID, sum)
	resp, _, errs := gorequest.New().Get(url).End()
	validStatusCode := resp.StatusCode >= 200 && resp.StatusCode < 300
	if len(errs) != 0 || validStatusCode {
		err := errors.Wrap(&multierror.Error{Errors: errs}, "fail to callback runcpa revenue share")
		logrus.WithFields(logrus.Fields{
			"event":            models.LogEventCallbackRuncpaRevenueShare,
			"url":              url,
			"http_status_code": resp.StatusCode,
			"error":            err.Error(),
		}).Warn("fail to callback runcpa revenue share")
	}
}
