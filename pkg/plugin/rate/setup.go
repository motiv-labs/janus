package rate

import (
	"net/http"

	"github.com/hellofresh/janus/pkg/errors"
	"github.com/hellofresh/janus/pkg/plugin"
	"github.com/hellofresh/janus/pkg/proxy"
	"github.com/hellofresh/janus/pkg/store"
	"github.com/ulule/limiter"
)

var (
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
	Limit  string `json:"limit"`
	Policy string `json:"policy"`
}

func init() {
	plugin.RegisterPlugin("rate_limit", plugin.Plugin{
		Action: setupRateLimit,
	})
}

func setupRateLimit(route *proxy.Route, p plugin.Params) error {
	var config Config
	err := plugin.Decode(p.Config, &config)
	if err != nil {
		return err
	}

	rate, err := limiter.NewRateFromFormatted(config.Limit)
	if err != nil {
		return err
	}

	limiterStore, err := getLimiterStore(p.Storage, config.Policy, route.Proxy.ListenPath)
	if err != nil {
		return err
	}

	limiterInstance := limiter.NewLimiter(limiterStore, rate)
	route.AddInbound(NewRateLimitLogger(limiterInstance, p.StatsClient))
	route.AddInbound(limiter.NewHTTPMiddleware(limiterInstance).Handler)

	return nil
}

func getLimiterStore(storage store.Store, policy string, prefix string) (limiter.Store, error) {
	if prefix == "" {
		prefix = DefaultPrefix
	}

	switch policy {
	case "redis":
		redisStorage, ok := storage.(*store.RedisStore)
		if !ok {
			return nil, ErrInvalidStorage
		}

		return limiter.NewRedisStoreWithOptions(redisStorage.Pool, limiter.StoreOptions{
			Prefix:   prefix,
			MaxRetry: limiter.DefaultMaxRetry,
		})
	case "local":
		return limiter.NewMemoryStore(), nil
	default:
		return nil, ErrInvalidPolicy
	}
}
