package main

import (
	"gopkg.in/mgo.v2/bson"
	"time"
	"github.com/RangelReale/osin"
)

type APIDefinition struct {
	ID             bson.ObjectId               `bson:"_id,omitempty" json:"id,omitempty"`
	CreatedAt      time.Time                   `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time                   `bson:"updated_at" json:"updated_at"`
	Name           string                      `bson:"name" json:"name"`
	Slug           string                      `bson:"slug" json:"slug"`
	Active         bool                        `bson:"active" json:"active"`
	UseBasicAuth   bool                        `bson:"use_basic_auth" json:"use_basic_auth"`
	Domain         string                      `bson:"domain" json:"domain"`
	Proxy          Proxy                       `bson:"proxy" json:"proxy"`
	AllowedIPs     []string                    `mapstructure:"allowed_ips" bson:"allowed_ips" json:"allowed_ips"`
	CircuitBreaker CircuitBreakerMeta          `bson:"circuit_breakers" json:"circuit_breakers"`
	UseOauth2      bool                        `bson:"use_oauth2" json:"use_oauth2"`
	Oauth2Meta     Oauth2Meta                  `bson:"oauth_meta" json:"oauth_meta"`
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

type Oauth2Meta       struct {
	AllowedAccessTypes     []osin.AccessRequestType    `bson:"allowed_access_types" json:"allowed_access_types"`
	AllowedAuthorizeTypes  []osin.AuthorizeRequestType `bson:"allowed_authorize_types" json:"allowed_authorize_types"`
	AuthorizeLoginRedirect string                      `bson:"auth_login_redirect" json:"auth_login_redirect"`
}
