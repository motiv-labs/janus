package provider

import (
	"net/http"
	"sync"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/hellofresh/janus/pkg/config"
)

var providers *sync.Map

func init() {
	providers = new(sync.Map)
}

// Provider represents an auth provider
type Provider interface {
	Verifier
	Build(config config.Credentials) Provider
	GetClaims(httpClient *http.Client) (jwt.MapClaims, error)
}

// Register registers a new provider
func Register(providerName string, providerConstructor Provider) {
	providers.Store(providerName, providerConstructor)
}

// GetProviders returns the list of registered providers
func GetProviders() *sync.Map {
	return providers
}

// Factory represents a factory of providers
type Factory struct{}

// Build builds one provider based on the auth configuration
func (f *Factory) Build(providerName string, config config.Credentials) Provider {
	provider, ok := providers.Load(providerName)
	if !ok {
		provider, _ = providers.Load("basic")
	}

	p := provider.(Provider)
	return p.Build(config)
}
