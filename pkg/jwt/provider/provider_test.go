package provider

import (
	"net/http"
	"testing"

	"github.com/hellofresh/janus/pkg/config"
	"github.com/stretchr/testify/assert"
)

type mockProvider struct{}

func (p *mockProvider) Build(config config.Credentials) Provider {
	return &mockProvider{}
}
func (p *mockProvider) Verify(r *http.Request) (bool, error) {
	return true, nil
}

type defaultProvider struct{}

func (p *defaultProvider) Build(config config.Credentials) Provider {
	return &defaultProvider{}
}
func (p *defaultProvider) Verify(r *http.Request) (bool, error) {
	return true, nil
}

func TestProviders(t *testing.T) {
	tests := []struct {
		scenario string
		function func(*testing.T, *Factory)
	}{
		{
			scenario: "it should build providers properly",
			function: testFactoryCanBuildProvider,
		},
		{
			scenario: "when given a wrong provider, it should get the default",
			function: testFactoryCantFindProvider,
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			t.Parallel()
			Register("test", &mockProvider{})
			Register("basic", &defaultProvider{})

			f := &Factory{}
			test.function(t, f)
			providers = make(map[string]Provider)
		})
	}
}

func testFactoryCanBuildProvider(t *testing.T, f *Factory) {
	p := f.Build("test", config.Credentials{})

	assert.Implements(t, (*Provider)(nil), p)
	assert.IsType(t, (*mockProvider)(nil), p)
}

func testFactoryCantFindProvider(t *testing.T, f *Factory) {
	p := f.Build("wrong", config.Credentials{})

	assert.Implements(t, (*Provider)(nil), p)
	assert.IsType(t, (*defaultProvider)(nil), p)
}

func testCountProvider(t *testing.T, f *Factory) {
	assert.Len(t, GetProviders(), 2)
}
