package plugin

import (
	"github.com/hellofresh/janus/pkg/api"
	"github.com/hellofresh/janus/pkg/middleware"
	"github.com/hellofresh/janus/pkg/router"
	"github.com/hellofresh/janus/pkg/store"
	"github.com/ulule/limiter"
)

// RateLimit represents the rate limit plugin
type RateLimit struct {
	storage store.Store
}

// NewRateLimit creates a new instance of HostMatcher
func NewRateLimit(storage store.Store) *RateLimit {
	return &RateLimit{storage}
}

// GetName retrieves the plugin's name
func (h *RateLimit) GetName() string {
	return "rate_limit"
}

// GetMiddlewares retrieves the plugin's middlewares
func (h *RateLimit) GetMiddlewares(config api.Config, referenceSpec *api.Spec) ([]router.Constructor, error) {
	limit := config["limit"].(string)
	rate, err := limiter.NewRateFromFormatted(limit)
	if err != nil {
		return nil, err
	}

	limiterStore, err := h.storage.ToLimiterStore(referenceSpec.Name)
	if err != nil {
		return nil, err
	}

	return []router.Constructor{
		limiter.NewHTTPMiddleware(limiter.NewLimiter(limiterStore, rate)).Handler,
		middleware.NewRateLimitLogger().Handler,
	}, nil
}
