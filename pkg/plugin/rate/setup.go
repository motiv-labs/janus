package rate

import (
	"net/http"
	"time"

	"github.com/go-redis/redis"
	"github.com/hellofresh/janus/pkg/api"
	"github.com/hellofresh/janus/pkg/errors"
	"github.com/hellofresh/janus/pkg/plugin"
	"github.com/hellofresh/janus/pkg/proxy"
	"github.com/hellofresh/stats-go/client"
	"github.com/ulule/limiter"
	"github.com/ulule/limiter/drivers/middleware/stdlib"
	smemory "github.com/ulule/limiter/drivers/store/memory"
	sredis "github.com/ulule/limiter/drivers/store/redis"
)

var (
	statsClient client.Client
	// ErrInvalidPolicy is used when an invalid policy was provided
	ErrInvalidPolicy = errors.New(http.StatusBadRequest, "policy is not supported")
	// ErrInvalidStorage is used when an invalid storage was provided
	ErrInvalidStorage = errors.New(http.StatusBadRequest, "the storage that you are using is not supported for this feature")
)

const (
	// DefaultPrefix is the default prefix to use for the key in the store.
	DefaultPrefix = "limiter"
)

// Config represents a rate limit config
type Config struct {
	Limit       string      `json:"limit"`
	Policy      string      `json:"policy"`
	RedisConfig redisConfig `json:"redis"`
}

type redisConfig struct {
	DSN    string `json:"dsn"`
	Prefix string `json:"prefix"`
}

func init() {
	plugin.RegisterEventHook(plugin.StartupEvent, onStartup)
	plugin.RegisterPlugin("rate_limit", plugin.Plugin{
		Action: setupRateLimit,
	})
}

func onStartup(event interface{}) error {
	e, ok := event.(plugin.OnStartup)
	if !ok {
		return errors.New(http.StatusInternalServerError, "Could not convert event to startup type")
	}

	statsClient = e.StatsClient
	return nil
}

func setupRateLimit(def *api.Definition, route *proxy.Route, rawConfig plugin.Config) error {
	var config Config
	err := plugin.Decode(rawConfig, &config)
	if err != nil {
		return err
	}

	rate, err := limiter.NewRateFromFormatted(config.Limit)
	if err != nil {
		return err
	}

	limiterStore, err := getLimiterStore(config.Policy, config.RedisConfig)
	if err != nil {
		return err
	}

	limiterInstance := limiter.New(limiterStore, rate)
	route.AddInbound(NewRateLimitLogger(limiterInstance, statsClient))
	route.AddInbound(stdlib.NewMiddleware(limiterInstance).Handler)

	return nil
}

func getLimiterStore(policy string, config redisConfig) (limiter.Store, error) {
	switch policy {
	case "redis":
		option, err := redis.ParseURL(config.DSN)
		if err != nil {
			return nil, err
		}
		option.PoolSize = 3
		option.IdleTimeout = 240 * time.Second
		client := redis.NewClient(option)

		if config.Prefix == "" {
			config.Prefix = DefaultPrefix
		}

		return sredis.NewStoreWithOptions(client, limiter.StoreOptions{
			Prefix:   config.Prefix,
			MaxRetry: limiter.DefaultMaxRetry,
		})
	case "local":
		return smemory.NewStore(), nil
	default:
		return nil, ErrInvalidPolicy
	}
}
