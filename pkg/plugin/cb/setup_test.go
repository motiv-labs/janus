package cb

import (
	"testing"

	"github.com/hellofresh/janus/pkg/config"
	"github.com/hellofresh/janus/pkg/plugin"
	"github.com/hellofresh/janus/pkg/proxy"
	"github.com/hellofresh/janus/pkg/router"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCbPlugin(t *testing.T) {
	t.Parallel()

	tests := []struct {
		scenario string
		function func(*testing.T)
	}{
		{
			scenario: "when the correct cb configuration is given",
			function: testSetupWithCorrectConfig,
		},
		{
			scenario: "when an incorrect cb configuration is given",
			function: testSetupWithIncorrectConfig,
		},
		{
			scenario: "when the plugin setup is successful",
			function: testSetupSuccess,
		},
		{
			scenario: "when the plugin admin startup is successful",
			function: testAdminStartupSuccess,
		},
		{
			scenario: "when the plugin startup is successful",
			function: testStartupSuccess,
		},
		{
			scenario: "when the plugin startup is not successful",
			function: testStartupNoSuccess,
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			test.function(t)
		})
	}
}

func testStartupNoSuccess(t *testing.T) {
	event2 := plugin.OnStartup{
		Register: proxy.NewRegister(proxy.WithRouter(router.NewChiRouter())),
		Config: &config.Specification{
			Stats: config.Stats{
				DSN: "statsd://statsd:8080",
			},
		},
	}
	err := onStartup(event2)
	require.Error(t, err)
}

func testStartupSuccess(t *testing.T) {
	event2 := plugin.OnStartup{
		Register: proxy.NewRegister(proxy.WithRouter(router.NewChiRouter())),
		Config:   &config.Specification{},
	}
	err := onStartup(event2)
	require.NoError(t, err)
}

func testAdminStartupSuccess(t *testing.T) {
	event1 := plugin.OnAdminAPIStartup{Router: router.NewChiRouter()}
	err := onAdminAPIStartup(event1)
	require.NoError(t, err)
}

func testSetupSuccess(t *testing.T) {
	def := proxy.NewRouterDefinition(proxy.NewDefinition())

	err := setupCB(def, make(plugin.Config))
	require.NoError(t, err)
}

func testSetupWithCorrectConfig(t *testing.T) {
	rawConfig := map[string]interface{}{
		"timeout":                 1000,
		"max_concurrent_requests": 100,
		"error_percent_threshold": 25,
		"sleep_window":            1,
		"predicate":               "statusCode => 500",
	}

	result, err := validateConfig(rawConfig)
	assert.True(t, result)
	require.NoError(t, err)
}

func testSetupWithIncorrectConfig(t *testing.T) {
	rawConfig := map[string]interface{}{
		"timeout": "wrong",
	}

	result, err := validateConfig(rawConfig)
	assert.False(t, result)
	require.Error(t, err)
}
