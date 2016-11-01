package main

import (
	"sync"

	"github.com/solefaucet/btcwall-api/models"
)

var (
	serviceState      = models.ServiceStateAlive
	serviceStateMutex sync.RWMutex
)

func isServiceAlive() bool {
	serviceStateMutex.RLock()
	defer serviceStateMutex.RUnlock()
	return serviceState == models.ServiceStateAlive
}

func setServiceToDead() {
	serviceStateMutex.Lock()
	defer serviceStateMutex.Unlock()
	serviceState = models.ServiceStateDead
}
