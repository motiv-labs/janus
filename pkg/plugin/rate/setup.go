package rate

import (
	"net/http"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/hellofresh/janus/pkg/errors"
	"github.com/hellofresh/janus/pkg/plugin"
	"github.com/hellofresh/janus/pkg/proxy"
	stats "github.com/hellofresh/stats-go"
	"github.com/ulule/limiter"
)

var (
	statsClient stats.Client
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

func setupRateLimit(route *proxy.Route, rawConfig plugin.Config) error {
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

	limiterInstance := limiter.NewLimiter(limiterStore, rate)
	route.AddInbound(NewRateLimitLogger(limiterInstance, statsClient))
	route.AddInbound(limiter.NewHTTPMiddleware(limiterInstance).Handler)

	return nil
}

func getLimiterStore(policy string, config redisConfig) (limiter.Store, error) {
	switch policy {
	case "redis":
		pool := &redis.Pool{
			MaxIdle:     3,
			IdleTimeout: 240 * time.Second,
			Dial:        func() (redis.Conn, error) { return redis.DialURL(config.DSN) },
		}

		if config.Prefix == "" {
			config.Prefix = DefaultPrefix
		}

		return limiter.NewRedisStoreWithOptions(pool, limiter.StoreOptions{
			Prefix:   config.Prefix,
			MaxRetry: limiter.DefaultMaxRetry,
		})
	case "local":
		return limiter.NewMemoryStore(), nil
	default:
		return nil, ErrInvalidPolicy
	}
}
