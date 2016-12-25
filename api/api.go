package api

import (
	"time"

	"github.com/hellofresh/janus/cors"
	"github.com/hellofresh/janus/proxy"
	"gopkg.in/mgo.v2/bson"
)

// Spec Holds an api definition and basic options
type Spec struct {
	Definition
}

// Definition Represents an API that you want to proxy
type Definition struct {
	ID            bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty" valid:"required"`
	CreatedAt     time.Time     `bson:"created_at" json:"created_at" valid:"-"`
	UpdatedAt     time.Time     `bson:"updated_at" json:"updated_at" valid:"-"`
	Name          string        `bson:"name" json:"name" valid:"required"`
	Slug          string        `bson:"slug" json:"slug"`
	Active        bool          `bson:"active" json:"active"`
	UseBasicAuth  bool          `bson:"use_basic_auth" json:"use_basic_auth"`
	Domain        string        `bson:"domain" json:"domain"`
	Proxy         proxy.Proxy   `bson:"proxy" json:"proxy" valid:"required"`
	AllowedIPs    []string      `mapstructure:"allowed_ips" bson:"allowed_ips" json:"allowed_ips"`
	UseOauth2     bool          `bson:"use_oauth2" json:"use_oauth2"`
	OAuthServerID bson.ObjectId `bson:"oauth_server_id" json:"oauth_server_id"`
	RateLimit     RateLimitMeta `bson:"rate_limit" json:"rate_limit" valid:"required"`
	CorsMeta      cors.Meta     `bson:"cors_meta" json:"cors_meta" valid:"cors_meta"`
}

// NewDefinition creates a new API Definition with default values
func NewDefinition() *Definition {
	return &Definition{
		OAuthServerID: bson.NewObjectId(),
	}
}

// RateLimitMeta holds configuration for a rate limit
type RateLimitMeta struct {
	Enabled bool   `bson:"enabled" json:"enabled"`
	Limit   string `bson:"limit" json:"limit"`
}
