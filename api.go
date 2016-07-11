package main

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type APIDefinition struct {
	ID             bson.ObjectId               `bson:"_id,omitempty" json:"id,omitempty"`
	CreatedAt      time.Time                   `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time                   `bson:"updated_at" json:"updated_at"`
	Name           string                      `bson:"name" json:"name"`
	Slug           string                      `bson:"slug" json:"slug"`
	Active         bool                        `bson:"active" json:"active"`
	Auth           Auth                        `bson:"auth" json:"auth"`
	UseBasicAuth   bool                        `bson:"use_basic_auth" json:"use_basic_auth"`
	Domain         string                      `bson:"domain" json:"domain"`
	Proxy          Proxy                       `bson:"proxy" json:"proxy"`
	AllowedIPs     []string                    `mapstructure:"allowed_ips" bson:"allowed_ips" json:"allowed_ips"`
	CircuitBreaker CircuitBreakerMeta        `bson:"circuit_breakers" json:"circuit_breakers"`
}

type Auth struct {
	UseParam       bool   `mapstructure:"use_param" bson:"use_param" json:"use_param"`
	ParamName      string `mapstructure:"param_name" bson:"param_name" json:"param_name"`
	UseCookie      bool   `mapstructure:"use_cookie" bson:"use_cookie" json:"use_cookie"`
	CookieName     string `mapstructure:"cookie_name" bson:"cookie_name" json:"cookie_name"`
	AuthHeaderName string `mapstructure:"auth_header_name" bson:"auth_header_name" json:"auth_header_name"`
}

type Proxy struct {
	PreserveHostHeader          bool                          `bson:"preserve_host_header" json:"preserve_host_header"`
	ListenPath                  string                        `bson:"listen_path" json:"listen_path"`
	TargetURL                   string                        `bson:"target_url" json:"target_url"`
	StripListenPath             bool                          `bson:"strip_listen_path" json:"strip_listen_path"`
	EnableLoadBalancing         bool                          `bson:"enable_load_balancing" json:"enable_load_balancing"`
	TargetList                  []string                      `bson:"target_list" json:"target_list"`
	CheckHostAgainstUptimeTests bool                          `bson:"check_host_against_uptime_tests" json:"check_host_against_uptime_tests"`
	//ServiceDiscovery            ServiceDiscoveryConfiguration `bson:"service_discovery" json:"service_discovery"`
}

type CircuitBreakerMeta struct {
	ThresholdPercent     float64 `bson:"threshold_percent" json:"threshold_percent"`
	Samples              int64   `bson:"samples" json:"samples"`
	ReturnToServiceAfter int     `bson:"return_to_service_after" json:"return_to_service_after"`
}
