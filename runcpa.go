package main

import "github.com/solefaucet/btcwall-api/service/runcpa"

var runcpaNotifier runcpa.Notifier

func initializeRuncpaNotifier() {
	runcpaNotifier = runcpa.New(config.Runcpa.BaseRegistrationCallbackURL, config.Runcpa.BaseRevenueShareCallbackURL)
}
