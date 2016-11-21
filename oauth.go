package janus

import (
	"time"

	"github.com/RangelReale/osin"
	"gopkg.in/mgo.v2/bson"
)

// OAuthSpec Holds an oauth definition and basic options
type OAuthSpec struct {
	*OAuth
	OAuthManager *OAuthManager
}

// OAuth holds the configuration for oauth proxies
type OAuth struct {
	ID                     bson.ObjectId               `bson:"_id,omitempty" json:"id,omitempty" valid:"required"`
	Name                   string                      `bson:"name" json:"name" valid:"required"`
	CreatedAt              time.Time                   `bson:"created_at" json:"created_at" valid:"-"`
	UpdatedAt              time.Time                   `bson:"updated_at" json:"updated_at" valid:"-"`
	OauthEndpoints         OauthEndpoints              `bson:"oauth_endpoints" json:"oauth_endpoints"`
	OauthClientEndpoints   OauthClientEndpoints        `bson:"oauth_client_endpoints" json:"oauth_client_endpoints"`
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
