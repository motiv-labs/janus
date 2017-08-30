package provider

import (
	"github.com/hellofresh/janus/pkg/config"
)

var providers map[string]Provider

// Provider represents an auth provider
type Provider interface {
	Verifier
	Build(config config.Credentials) Provider
}

func init() {
	providers = make(map[string]Provider)
}

// Register registers a new provider
func Register(providerName string, providerConstructor Provider) {
	providers[providerName] = providerConstructor
}

// GetProviders returns the list of registered providers
func GetProviders() map[string]Provider {
	return providers
}

// Factory represents a factory of providers
type Factory struct{}

// Build builds one provider based on the auth configuration
func (f *Factory) Build(providerName string, config config.Credentials) Provider {
	provider, ok := providers[providerName]
	if !ok {
		provider = providers["basic"]
	}
	return provider.Build(config)
}
