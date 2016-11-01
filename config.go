package main

import (
	"os"

	"github.com/solefaucet/btcwall-api/models"
	"github.com/spf13/viper"
)

var config models.Configuration

func initializeConfiguration() {
	filename := os.Getenv("BTCWALL_API_CONF")
	if filename == "" {
		filename = "./conf.yml"
	}

	viper.SetConfigType("yaml")
	viper.SetConfigFile(filename)

	must(nil, viper.ReadInConfig())
	must(nil, viper.Unmarshal(&config))
	must(nil, config.Validate())
}
