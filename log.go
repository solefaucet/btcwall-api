package main

import (
	"os"

	"github.com/Sirupsen/logrus"
	graylog "github.com/on99/logrus-graylog-hook"
	"github.com/solefaucet/btcwall-api/models"
)

var logger logging

func initializeLogging() {
	l := must(logrus.ParseLevel(config.Log.Level)).(logrus.Level)
	logrus.SetLevel(l)
	logrus.SetOutput(os.Stdout)

	if config.Log.GraylogHook != nil {
		c := graylog.NewConfig()
		c.Facility = config.Log.GraylogHook.Facility
		c.HealthCheckInterval = config.Log.GraylogHook.HealthCheckInterval
		c.StaticMeta = map[string]interface{}{
			"go_version": goVersion,
			"build_time": buildTime,
			"git_commit": gitCommit,
		}

		if logger.graylogHook == nil {
			hook := graylog.New(c)
			hook.StartHealthCheck()
			logrus.AddHook(hook)

			logger.graylogHook = hook
		}

		logger.graylogHook.SetNodeConfigs(graylogNodesToGraylogNodeConfigs(config.Log.GraylogHook.Nodes)...)
	}
}

type logging struct {
	graylogHook *graylog.Hook
}

func graylogNodesToGraylogNodeConfigs(nodes []models.GraylogNode) []graylog.NodeConfig {
	var ns []graylog.NodeConfig
	for _, node := range nodes {
		ns = append(ns, graylog.NodeConfig{
			UDPAddress:     node.UDPAddress,
			HealthCheckURL: node.HealthCheckURL,
			Weight:         node.Weight,
		})
	}
	return ns
}
