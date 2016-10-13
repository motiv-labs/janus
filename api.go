package main

import (
	"time"

	"github.com/RangelReale/osin"
	"gopkg.in/mgo.v2/bson"
)

// APISpec Holds an api definition and basic options
type APISpec struct {
	APIDefinition
	OAuthManager *OAuthManager
}

// APIDefinition Represents an API that you want to proxy
type APIDefinition struct {
	ID           bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty" valid:"required"`
	CreatedAt    time.Time     `bson:"created_at" json:"created_at" valid:"-"`
	UpdatedAt    time.Time     `bson:"updated_at" json:"updated_at" valid:"-"`
	Name         string        `bson:"name" json:"name" valid:"required"`
	Slug         string        `bson:"slug" json:"slug"`
	Active       bool          `bson:"active" json:"active"`
	UseBasicAuth bool          `bson:"use_basic_auth" json:"use_basic_auth"`
	Domain       string        `bson:"domain" json:"domain"`
	Proxy        Proxy         `bson:"proxy" json:"proxy" valid:"required"`
	AllowedIPs   []string      `mapstructure:"allowed_ips" bson:"allowed_ips" json:"allowed_ips"`
	UseOauth2    bool          `bson:"use_oauth2" json:"use_oauth2"`
	Oauth2Meta   Oauth2Meta    `bson:"oauth_meta" json:"oauth_meta" valid:"required"`
	RateLimit    RateLimitMeta `bson:"rate_limit" json:"rate_limit" valid:"required"`
	CorsMeta     CorsMeta      `bson:"cors_meta" json:"cors_meta" valid:"cors_meta"`
}

// Proxy defines proxy rules for a route
type Proxy struct {
	PreserveHostHeader          bool     `bson:"preserve_host_header" json:"preserve_host_header"`
	ListenPath                  string   `bson:"listen_path" json:"listen_path" valid:"required"`
	TargetURL                   string   `bson:"target_url" json:"target_url" valid:"url,required"`
	StripListenPath             bool     `bson:"strip_listen_path" json:"strip_listen_path"`
	EnableLoadBalancing         bool     `bson:"enable_load_balancing" json:"enable_load_balancing"`
	TargetList                  []string `bson:"target_list" json:"target_list"`
	CheckHostAgainstUptimeTests bool     `bson:"check_host_against_uptime_tests" json:"check_host_against_uptime_tests"`
	Methods                     []string `bson:"methods" json:"methods"`
}

// CorsMeta defines config for CORS routes
type CorsMeta struct {
	Domains        []string `mapstructure:"domains" bson:"domains" json:"domains"`
	Methods        []string `mapstructure:"methods" bson:"methods" json:"methods"`
	RequestHeaders []string `mapstructure:"request_headers" bson:"request_headers" json:"request_headers"`
	ExposedHeaders []string `mapstructure:"exposed_headers" bson:"exposed_headers" json:"exposed_headers"`
	Enabled        bool     `bson:"enabled" json:"enabled"`
}

// Oauth2Meta holds the configuration for oauth proxies
type Oauth2Meta struct {
	OauthEndpoints         `bson:"oauth_endpoints" json:"oauth_endpoints"`
	OauthClientEndpoints   `bson:"oauth_client_endpoints" json:"oauth_client_endpoints"`
	AllowedAccessTypes     []osin.AccessRequestType    `bson:"allowed_access_types" json:"allowed_access_types"`
	AllowedAuthorizeTypes  []osin.AuthorizeRequestType `bson:"allowed_authorize_types" json:"allowed_authorize_types"`
	AuthorizeLoginRedirect string                      `bson:"auth_login_redirect" json:"auth_login_redirect"`
}

// OauthEndpoints defines the oauth endpoints that wil be proxied
type OauthEndpoints struct {
	Authorize Proxy `bson:"authorize" json:"authorize"`
	Token     Proxy `bson:"token" json:"token"`
	Info      Proxy `bson:"info" json:"info"`
}

// OauthClientEndpoints defines the oauth client endpoints that wil be proxied
type OauthClientEndpoints struct {
	Create Proxy `bson:"create" json:"create"`
	Remove Proxy `bson:"remove" json:"remove"`
}

// RateLimitMeta holds configuration for a rate limit
type RateLimitMeta struct {
	Enabled bool  `bson:"enabled" json:"enabled"`
	Limit   int64 `bson:"limit" json:"limit"`
}
