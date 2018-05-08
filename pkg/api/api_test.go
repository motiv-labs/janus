package api_test

import (
	"testing"

	"github.com/hellofresh/janus/pkg/api"
	"github.com/hellofresh/janus/pkg/proxy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewInstanceOfDefinition(t *testing.T) {
	instance := api.NewDefinition()

	assert.IsType(t, &api.Definition{}, instance)
	assert.True(t, instance.Active)
}

func TestSuccessfulValidation(t *testing.T) {
	instance := api.NewDefinition()
	instance.Name = "Test"
	instance.Proxy.ListenPath = "/"
	instance.Proxy.Upstreams = &proxy.Upstreams{
		Balancing: "roundrobin",
		Targets: []*proxy.Target{
			{Target: "http:/example.com"},
		},
	}

	isValid, err := instance.Validate()
	require.NoError(t, err)
	assert.True(t, isValid)
}

func TestFailedValidation(t *testing.T) {
	instance := api.NewDefinition()
	isValid, err := instance.Validate()

	assert.Error(t, err)
	assert.False(t, isValid)
}

func TestNameValidation(t *testing.T) {
	instanceSimple := api.NewDefinition()
	instanceSimple.Name = "simple"
	instanceSimple.Proxy.ListenPath = "/"
	isValid, err := instanceSimple.Validate()

	require.NoError(t, err)
	require.True(t, isValid)

	instanceDash := api.NewDefinition()
	instanceDash.Name = "with-dash-and-123"
	instanceDash.Proxy.ListenPath = "/"
	isValid, err = instanceDash.Validate()

	require.NoError(t, err)
	require.True(t, isValid)

	instanceBadSymbol := api.NewDefinition()
	instanceBadSymbol.Name = "test~"
	instanceBadSymbol.Proxy.ListenPath = "/"
	isValid, err = instanceBadSymbol.Validate()

	require.Error(t, err)
	require.False(t, isValid)
}
