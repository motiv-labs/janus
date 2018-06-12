package proxy

import (
	"time"

	"github.com/hellofresh/janus/pkg/router"
	"github.com/hellofresh/stats-go/client"
)

// RegisterOption represents the register options
type RegisterOption func(*Register)

// WithRouter sets the router
func WithRouter(router router.Router) RegisterOption {
	return func(r *Register) {
		r.router = router
	}
}

// WithFlushInterval sets the Flush interval for copying upgraded connections
func WithFlushInterval(d time.Duration) RegisterOption {
	return func(r *Register) {
		r.flushInterval = d
	}
}

// WithCloseIdleConnsPeriod sets the time period of how often the idle connections are
// forcibly closed
func WithCloseIdleConnsPeriod(d time.Duration) RegisterOption {
	return func(r *Register) {
		r.closeIdleConnsPeriod = d
	}
}

// WithIdleConnectionsPerHost sets idle connections per host option
func WithIdleConnectionsPerHost(value int) RegisterOption {
	return func(r *Register) {
		r.idleConnectionsPerHost = value
	}
}

// WithStatsClient sets stats client instance for proxy
func WithStatsClient(statsClient client.Client) RegisterOption {
	return func(r *Register) {
		r.statsClient = statsClient
	}
}
