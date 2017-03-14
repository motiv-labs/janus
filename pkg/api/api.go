package api

import (
	"github.com/hellofresh/janus/pkg/cors"
	"github.com/hellofresh/janus/pkg/oauth"
	"github.com/hellofresh/janus/pkg/proxy"
)

// Spec Holds an api definition and basic options
type Spec struct {
	*Definition
	Manager oauth.Manager
}

// Definition Represents an API that you want to proxy
type Definition struct {
	Name            string            `bson:"name" json:"name" valid:"required"`
	Active          bool              `bson:"active" json:"active"`
	Proxy           *proxy.Definition `bson:"proxy" json:"proxy" valid:"required"`
	AllowedIPs      []string          `mapstructure:"allowed_ips" bson:"allowed_ips" json:"allowed_ips"`
	UseBasicAuth    bool              `bson:"use_basic_auth" json:"use_basic_auth"`
	UseOauth2       bool              `bson:"use_oauth2" json:"use_oauth2"`
	OAuthServerSlug string            `bson:"oauth_server_slug" json:"oauth_server_slug"`
	RateLimit       RateLimitMeta     `bson:"rate_limit" json:"rate_limit" valid:"required"`
	CorsMeta        cors.Meta         `bson:"cors_meta" json:"cors_meta" valid:"cors_meta"`
	UseCompression  bool              `bson:"use_compression" json:"use_compression"`
}

// NewDefinition creates a new API Definition with default values
func NewDefinition() *Definition {
	return &Definition{
		UseCompression: true,
	}
}

// RateLimitMeta holds configuration for a rate limit
type RateLimitMeta struct {
	Enabled bool   `bson:"enabled" json:"enabled"`
	Limit   string `bson:"limit" json:"limit"`
}
