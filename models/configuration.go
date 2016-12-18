package models

import (
	"time"

	validator "gopkg.in/go-playground/validator.v8"
)

// Configuration model
type Configuration struct {
	HTTP struct {
		Address string `mapstructure:"address" validate:"required,tcp_addr"`
		Mode    string `mapstructure:"mode" validate:"required,eq=release|eq=test|eq=debug"`
	} `mapstructure:"http" validate:"required"`
	Log struct {
		Level       string       `mapstructure:"level" validate:"required,eq=debug|eq=info|eq=warn|eq=error|eq=fatal|eq=panic"`
		GraylogHook *GraylogHook `mapstructure:"graylog_hook" validate:"omitempty"`
	} `mapstructure:"log" validate:"required"`
	Geo struct {
		Filename string `mapstructure:"filename" validate:"required"`
	} `mapstructure:"geo" validate:"required"`
	ProxyDetection struct {
		Threshold int64 `mapstructure:"threshold" validate:"required,min=0,max=5"`
	} `mapstructure:"proxy_detection" validate:"required"`
	RPC struct {
		Address               string `mapstructure:"address" validate:"required,url"`
		MaxIdleConnsPerHost   int    `mapstructure:"max_idle_conns_per_host" validate:"required,min=1"`
		MaxConcurrentRequests int    `mapstructure:"max_concurrent_requests" validate:"required,min=1"`
	} `mapstructure:"rpc" validate:"required"`
	Offerwalls struct {
		Adgate struct {
			WhitelistIPs []string `mapstructure:"whitelist_ips" validate:"omitempty,dive,ip"`
		} `mapstructure:"adgate"`
		Adscend struct {
			WhitelistIPs []string `mapstructure:"whitelist_ips" validate:"omitempty,dive,ip"`
		} `mapstructure:"adscend"`
		Kiwiwall struct {
			WhitelistIPs []string `mapstructure:"whitelist_ips" validate:"omitempty,dive,ip"`
			SecretKey    string   `mapstructure:"secret_key"`
		} `mapstructure:"kiwiwall"`
		Personaly struct {
			WhitelistIPs []string `mapstructure:"whitelist_ips" validate:"omitempty,dive,ip"`
			AppHash      string   `mapstructure:"app_hash"`
			SecretKey    string   `mapstructure:"secret_key"`
		} `mapstructure:"personaly"`
		PointClickTrack struct {
			WhitelistIPs []string `mapstructure:"whitelist_ips" validate:"omitempty,dive,ip"`
		} `mapstructure:"point_click_track"`
		Ptcwall struct {
			WhitelistIPs     []string `mapstructure:"whitelist_ips" validate:"omitempty,dive,ip"`
			PostbackPassword string   `mapstructure:"postback_password"`
		} `mapstructure:"ptcwall"`
		Wannads struct {
			SecretKey string `mapstructure:"secret_key"`
		} `mapstructure:"wannads"`
		Peanut struct {
			WhitelistIPs   []string `mapstructure:"whitelist_ips" validate:"omitempty,dive,ip"`
			ApplicationKey string   `mapstructure:"application_key"`
			TransactionKey string   `mapstructure:"transaction_key"`
		} `mapstructure:"peanut"`
	} `mapstructure:"offerwalls"`
	Runcpa struct {
		BaseRegistrationCallbackURL string `mapstructure:"base_registration_callback_url" validate:"required,url"`
		BaseRevenueShareCallbackURL string `mapstructure:"base_revenue_share_callback_url" validate:"required,url"`
	} `mapstructure:"runcpa" validate:"required"`
}

// GraylogHook model
type GraylogHook struct {
	Facility            string        `mapstructure:"facility" validate:"required"`
	HealthCheckInterval time.Duration `mapstructure:"health_check_interval" validate:"required"`
	Nodes               []GraylogNode `mapstructure:"nodes" validate:"omitempty,dive"`
}

// GraylogNode model
type GraylogNode struct {
	UDPAddress     string `mapstructure:"udp_address" validate:"udp_addr"`
	HealthCheckURL string `mapstructure:"health_check_url" validate:"url"`
	Weight         int    `mapstructure:"weight" validate:"required,min=1"`
}

// Validate validates configuration
func (c Configuration) Validate() error {
	validate := validator.New(&validator.Config{TagName: "validate"})
	return validate.Struct(c)
}

// Copy config
func (c *Configuration) Copy(config Configuration) {
	c.HTTP = config.HTTP
	c.Log = config.Log
	c.Geo = config.Geo
	c.RPC = config.RPC
	c.Offerwalls = config.Offerwalls
	c.Runcpa = config.Runcpa
}
