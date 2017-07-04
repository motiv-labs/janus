package rate

import (
	"testing"

	"github.com/garyburd/redigo/redis"
	"github.com/hellofresh/janus/pkg/plugin"
	"github.com/hellofresh/janus/pkg/proxy"
	"github.com/hellofresh/janus/pkg/store"
	stats "github.com/hellofresh/stats-go"
	"github.com/rafaeljusto/redigomock"
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

	statsClient, _ := stats.NewClient("memory://", "")
	route := proxy.NewRoute(&proxy.Definition{})
	err := setupRateLimit(route, plugin.Params{
		Config:      rawConfig,
		Storage:     store.NewInMemoryStore(),
		StatsClient: statsClient,
	})

	assert.NoError(t, err)
	assert.Len(t, route.Inbound, 2)
}

func TestRateLimitPluginRedisPolicyWithInvalidStorage(t *testing.T) {
	rawConfig := map[string]interface{}{
		"limit":  "10-S",
		"policy": "redis",
	}

	statsClient, _ := stats.NewClient("memory://", "")
	route := proxy.NewRoute(&proxy.Definition{})
	err := setupRateLimit(route, plugin.Params{
		Config:      rawConfig,
		Storage:     store.NewInMemoryStore(),
		StatsClient: statsClient,
	})

	assert.Error(t, err)
}

func TestRateLimitPluginRedisPolicy(t *testing.T) {
	rawConfig := map[string]interface{}{
		"limit":  "10-S",
		"policy": "redis",
	}

	pool := redis.NewPool(func() (redis.Conn, error) {
		return redigomock.NewConn(), nil
	}, 0)
	storage, err := store.NewRedisStore(pool, "")
	assert.NoError(t, err)

	statsClient, _ := stats.NewClient("memory://", "")
	route := proxy.NewRoute(&proxy.Definition{})
	err = setupRateLimit(route, plugin.Params{
		Config:      rawConfig,
		Storage:     storage,
		StatsClient: statsClient,
	})

	assert.Error(t, err)
}

func TestRateLimitPluginInvalidPolicy(t *testing.T) {
	rawConfig := map[string]interface{}{
		"limit":  "10-S",
		"policy": "wrong",
	}

	statsClient, _ := stats.NewClient("memory://", "")
	route := proxy.NewRoute(&proxy.Definition{})
	err := setupRateLimit(route, plugin.Params{
		Config:      rawConfig,
		Storage:     store.NewInMemoryStore(),
		StatsClient: statsClient,
	})

	assert.Error(t, err)
}
