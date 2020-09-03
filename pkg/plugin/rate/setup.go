package rate

import (
	"net/http"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/go-redis/redis/v7"
	"github.com/hellofresh/stats-go/client"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/middleware/stdlib"
	storeMemory "github.com/ulule/limiter/v3/drivers/store/memory"
	storeRedis "github.com/ulule/limiter/v3/drivers/store/redis"

	"github.com/hellofresh/janus/pkg/errors"
	"github.com/hellofresh/janus/pkg/plugin"
	"github.com/hellofresh/janus/pkg/proxy"
)

var (
	statsClient client.Client
	// ErrInvalidPolicy is used when an invalid policy was provided
	ErrInvalidPolicy = errors.New(http.StatusBadRequest, "policy is not supported")
)

const (
	// DefaultPrefix is the default prefix to use for the key in the store.
	DefaultPrefix = "limiter"
)

// Config represents a rate limit config
type Config struct {
	Limit               string      `json:"limit"`
	Policy              string      `json:"policy"`
	RedisConfig         redisConfig `json:"redis"`
	TrustForwardHeaders bool        `json:"trust_forward_headers"`
}

type redisConfig struct {
	DSN    string `json:"dsn"`
	Prefix string `json:"prefix"`
}

func init() {
	plugin.RegisterEventHook(plugin.StartupEvent, onStartup)
	plugin.RegisterPlugin("rate_limit", plugin.Plugin{
		Action:   setupRateLimit,
		Validate: validateConfig,
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

func validateConfig(rawConfig plugin.Config) (bool, error) {
	var config Config
	err := plugin.Decode(rawConfig, &config)
	if err != nil {
		return false, err
	}

	return govalidator.ValidateStruct(config)
}

func setupRateLimit(def *proxy.RouterDefinition, rawConfig plugin.Config) error {
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

	limiterInstance := limiter.New(limiterStore, rate, limiter.WithTrustForwardHeader(config.TrustForwardHeaders))
	def.AddMiddleware(NewRateLimitLogger(limiterInstance, statsClient, config.TrustForwardHeaders))
	def.AddMiddleware(stdlib.NewMiddleware(limiterInstance).Handler)

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
		redisClient := redis.NewClient(option)

		if config.Prefix == "" {
			config.Prefix = DefaultPrefix
		}

		return storeRedis.NewStoreWithOptions(redisClient, limiter.StoreOptions{
			Prefix:   config.Prefix,
			MaxRetry: limiter.DefaultMaxRetry,
		})

	case "local":
		return storeMemory.NewStore(), nil

	default:
		return nil, ErrInvalidPolicy
	}
}
