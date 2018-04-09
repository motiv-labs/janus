package oauth2

import (
	"sync"

	"github.com/Knetic/govaluate"
	"github.com/hellofresh/janus/pkg/jwt"
	"github.com/hellofresh/janus/pkg/proxy"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

// AccessRequestType is the type for OAuth param `grant_type`
type AccessRequestType string

// AuthorizeRequestType is the type for OAuth param `response_type`
type AuthorizeRequestType string

// Spec Holds an api definition and basic options
type Spec struct {
	*OAuth
	Manager Manager
}

// OAuth holds the configuration for oauth proxies
type OAuth struct {
	Name                   string                 `bson:"name" json:"name" valid:"required"`
	Endpoints              Endpoints              `bson:"oauth_endpoints" json:"oauth_endpoints" mapstructure:"oauth_endpoints"`
	ClientEndpoints        ClientEndpoints        `bson:"oauth_client_endpoints" json:"oauth_client_endpoints" mapstructure:"oauth_client_endpoints"`
	AllowedAccessTypes     []AccessRequestType    `bson:"allowed_access_types" json:"allowed_access_types" mapstructure:"allowed_access_types" `
	AllowedAuthorizeTypes  []AuthorizeRequestType `bson:"allowed_authorize_types" json:"allowed_authorize_types" mapstructure:"allowed_authorize_types"`
	AuthorizeLoginRedirect string                 `bson:"auth_login_redirect" json:"auth_login_redirect" mapstructure:"auth_login_redirect"`
	Secrets                map[string]string      `bson:"secrets" json:"secrets"`
	CorsMeta               corsMeta               `bson:"cors_meta" json:"cors_meta" mapstructure:"cors_meta"`
	RateLimit              rateLimitMeta          `bson:"rate_limit" json:"rate_limit"`
	TokenStrategy          TokenStrategy          `bson:"token_strategy" json:"token_strategy" mapstructure:"token_strategy"`
	AccessRules            []*AccessRule          `bson:"access_rules" json:"access_rules"`
}

// Endpoints defines the oauth endpoints that wil be proxied
type Endpoints struct {
	Authorize  *proxy.Definition `bson:"authorize" json:"authorize"`
	Token      *proxy.Definition `bson:"token" json:"token"`
	Introspect *proxy.Definition `bson:"introspect" json:"introspect"`
	Revoke     *proxy.Definition `bson:"revoke" json:"revoke"`
}

// ClientEndpoints defines the oauth client endpoints that wil be proxied
type ClientEndpoints struct {
	Create *proxy.Definition `bson:"create" json:"create"`
	Remove *proxy.Definition `bson:"remove" json:"remove"`
}

type rateLimitMeta struct {
	Limit   string `bson:"limit" json:"limit"`
	Enabled bool   `bson:"enabled" json:"enabled"`
}

type corsMeta struct {
	Domains        []string `mapstructure:"domains" bson:"domains" json:"domains"`
	Methods        []string `mapstructure:"methods" bson:"methods" json:"methods"`
	RequestHeaders []string `mapstructure:"request_headers" bson:"request_headers" json:"request_headers"`
	ExposedHeaders []string `mapstructure:"exposed_headers" bson:"exposed_headers" json:"exposed_headers"`
	Enabled        bool     `bson:"enabled" json:"enabled"`
}

// IntrospectionSettings represents the settings for introspection
type IntrospectionSettings struct {
	UseCustomHeader bool   `mapstructure:"use_custom_header" bson:"use_custom_header" json:"use_custom_header"`
	HeaderName      string `mapstructure:"header_name" bson:"header_name" json:"header_name"`
	UseAuthHeader   bool   `mapstructure:"use_auth_header" bson:"use_auth_header" json:"use_auth_header"`
	AuthHeaderType  string `mapstructure:"auth_header_type" bson:"auth_header_type" json:"auth_header_type"`
	UseBody         bool   `mapstructure:"use_body" bson:"use_body" json:"use_body"`
	ParamName       string `mapstructure:"param_name" bson:"param_name" json:"param_name"`
}

// TokenStrategy defines the token strategy fields
type TokenStrategy struct {
	Name     string      `bson:"name" json:"name"`
	Settings interface{} `bson:"settings" json:"settings"`
	// TODO: this should become part of the settings, but for "jwt" strategy we expect array of signing methods
	// at the moment, so this will be BC-breaking change. In the next major version we need to turn settings
	// into object/dictionary and make leeway one of the settings.
	Leeway int64 `bson:"leeway" json:"leeway"`
}

// NewOAuth creates a new instance of OAuth
func NewOAuth() *OAuth {
	return &OAuth{
		Secrets: make(map[string]string),
		Endpoints: Endpoints{
			Authorize:  proxy.NewDefinition(),
			Introspect: proxy.NewDefinition(),
			Revoke:     proxy.NewDefinition(),
			Token:      proxy.NewDefinition(),
		},
	}
}

// GetIntrospectionSettings returns the settings for introspection
func (t TokenStrategy) GetIntrospectionSettings() (*IntrospectionSettings, error) {
	var settings *IntrospectionSettings
	err := mapstructure.Decode(t.Settings, &settings)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode introspection settings")
	}
	return settings, nil
}

// GetJWTSigningMethods parses and returns chain of JWT signing methods for token signature validation.
// Supports fallback to legacy format with {"secret": "key"} as single signing method with HS256 alg.
func (t TokenStrategy) GetJWTSigningMethods() ([]jwt.SigningMethod, error) {
	var methods []jwt.SigningMethod
	err := mapstructure.Decode(t.Settings, &methods)
	// TODO: remove legacy format support in couple minor releases
	if err != nil {
		var legacy struct {
			Secret string `json:"secret"`
		}
		err = mapstructure.Decode(t.Settings, &legacy)
		if nil != err {
			return methods, err
		}
		if legacy.Secret == "" {
			return nil, ErrJWTSecretMissing
		}

		return []jwt.SigningMethod{{Alg: "HS256", Key: legacy.Secret}}, nil
	}
	return methods, err
}

// AccessRule represents a rule that will be applied to a JWT that could be revoked
type AccessRule struct {
	mu        sync.Mutex
	Predicate string `bson:"predicate" json:"predicate"`
	Action    string `bson:"action" json:"action"`
	parsed    bool
}

// IsAllowed checks if the rule is allowed to
func (r *AccessRule) IsAllowed(claims map[string]interface{}) (bool, error) {
	var err error

	if !r.parsed {
		matched, err := r.parse(claims)
		if err != nil {
			return false, err
		}

		if !matched {
			return true, nil
		}
	}

	return r.Action == "allow", err
}

func (r *AccessRule) parse(claims map[string]interface{}) (bool, error) {
	expression, err := govaluate.NewEvaluableExpression(r.Predicate)
	if err != nil {
		return false, errors.New("Could not create an expression with this predicate")
	}

	result, err := expression.Evaluate(claims)
	if err != nil {
		return false, errors.New("Cannot evaluate the expression")
	}

	r.mu.Lock()
	r.parsed = true
	r.mu.Unlock()

	return result.(bool), nil
}
