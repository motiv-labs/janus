package rate

import (
	"testing"

	"github.com/hellofresh/janus/pkg/plugin"
	"github.com/hellofresh/janus/pkg/proxy"
	"github.com/stretchr/testify/assert"
)

func TestRateLimitConfig(t *testing.T) {
	var config Config
	rawConfig := map[string]interface{}{
		"limit":  "10-S",
		"policy": "local",
	}

	err := plugin.Decode(rawConfig, &config)
	assert.NoError(t, err)

	assert.Equal(t, "10-S", config.Limit)
	assert.Equal(t, "local", config.Policy)
}

func TestInvalidRateLimitConfig(t *testing.T) {
	var config Config
	rawConfig := map[string]interface{}{
		"limit": []string{"wrong"},
	}

	err := plugin.Decode(rawConfig, &config)
	assert.Error(t, err)
}

func TestRateLimitPluginLocalPolicy(t *testing.T) {
	rawConfig := map[string]interface{}{
		"limit":  "10-S",
		"policy": "local",
	}

	route := proxy.NewRoute(&proxy.Definition{})
	err := setupRateLimit(route, rawConfig)

	assert.NoError(t, err)
	assert.Len(t, route.Inbound, 2)
}

func TestRateLimitPluginRedisPolicyWithInvalidStorage(t *testing.T) {
	rawConfig := map[string]interface{}{
		"limit":  "10-S",
		"policy": "redis",
	}

	route := proxy.NewRoute(&proxy.Definition{})
	err := setupRateLimit(route, rawConfig)

	assert.Error(t, err)
}

func TestRateLimitPluginInvalidPolicy(t *testing.T) {
	rawConfig := map[string]interface{}{
		"limit":  "10-S",
		"policy": "wrong",
	}

	route := proxy.NewRoute(&proxy.Definition{})
	err := setupRateLimit(route, rawConfig)

	assert.Error(t, err)
}
