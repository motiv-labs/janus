package plugin

import (
	"testing"

	"github.com/garyburd/redigo/redis"
	"github.com/hellofresh/janus/pkg/api"
	"github.com/hellofresh/janus/pkg/store"
	stats "github.com/hellofresh/stats-go"
	"github.com/rafaeljusto/redigomock"
	"github.com/stretchr/testify/assert"
)

func TestRateLimitConfig(t *testing.T) {
	var config rateLimitConfig
	rawConfig := map[string]interface{}{
		"limit":  "10-S",
		"policy": "local",
	}

	err := decode(rawConfig, &config)
	assert.NoError(t, err)

	assert.Equal(t, "10-S", config.Limit)
	assert.Equal(t, "local", config.Policy)
}

func TestInvalidRateLimitConfig(t *testing.T) {
	var config rateLimitConfig
	rawConfig := map[string]interface{}{
		"limit": []string{"wrong"},
	}

	err := decode(rawConfig, &config)
	assert.Error(t, err)
}

func TestRateLimitPluginGetName(t *testing.T) {
	statsClient, _ := stats.NewClient("memory://", "")
	plugin := NewRateLimit(store.NewInMemoryStore(), statsClient)

	assert.Equal(t, "rate_limit", plugin.GetName())
}

func TestRateLimitPluginLocalPolicy(t *testing.T) {
	rawConfig := map[string]interface{}{
		"limit":  "10-S",
		"policy": "local",
	}

	spec := &api.Spec{
		Definition: &api.Definition{
			Name: "API Name",
		},
	}

	statsClient, _ := stats.NewClient("memory://", "")
	plugin := NewRateLimit(store.NewInMemoryStore(), statsClient)
	middleware, err := plugin.GetMiddlewares(rawConfig, spec)

	assert.NoError(t, err)
	assert.Len(t, middleware, 2)
}

func TestRateLimitPluginRedisPolicyWithInvalidStorage(t *testing.T) {
	rawConfig := map[string]interface{}{
		"limit":  "10-S",
		"policy": "redis",
	}

	spec := &api.Spec{
		Definition: &api.Definition{
			Name: "API Name",
		},
	}

	statsClient, _ := stats.NewClient("memory://", "")
	plugin := NewRateLimit(store.NewInMemoryStore(), statsClient)
	_, err := plugin.GetMiddlewares(rawConfig, spec)

	assert.Error(t, err)
}

func TestRateLimitPluginRedisPolicy(t *testing.T) {
	rawConfig := map[string]interface{}{
		"limit":  "10-S",
		"policy": "redis",
	}

	spec := &api.Spec{
		Definition: &api.Definition{
			Name: "API Name",
		},
	}

	pool := redis.NewPool(func() (redis.Conn, error) {
		return redigomock.NewConn(), nil
	}, 0)
	storage, err := store.NewRedisStore(pool, "")
	assert.NoError(t, err)

	statsClient, _ := stats.NewClient("memory://", "")
	plugin := NewRateLimit(storage, statsClient)
	_, err = plugin.GetMiddlewares(rawConfig, spec)

	assert.Error(t, err)
}

func TestRateLimitPluginInvalidPolicy(t *testing.T) {
	rawConfig := map[string]interface{}{
		"limit":  "10-S",
		"policy": "wrong",
	}

	spec := &api.Spec{
		Definition: &api.Definition{
			Name: "API Name",
		},
	}

	statsClient, _ := stats.NewClient("memory://", "")
	plugin := NewRateLimit(store.NewInMemoryStore(), statsClient)
	_, err := plugin.GetMiddlewares(rawConfig, spec)

	assert.Error(t, err)
}
