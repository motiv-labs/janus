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

// WithIdleConnTimeout sets the maximum amount of time an idle
// (keep-alive) connection will remain idle before closing
// itself.
func WithIdleConnTimeout(d time.Duration) RegisterOption {
	return func(r *Register) {
		r.idleConnTimeout = d
	}
}

// WithIdleConnPurgeTicker purges idle connections on every interval if set
// this is done to prevent permanent keep-alive on connections with high ops
func WithIdleConnPurgeTicker(d time.Duration) RegisterOption {
	var ticker *time.Ticker

	if d != 0 {
		ticker = time.NewTicker(d)
	}

	return func(t *Register) {
		t.idleConnPurgeTicker = ticker
	}
}

// WithIsPublicEndpoint adds trace metadata from incoming requests
// as parent span if set to false
func WithIsPublicEndpoint(b bool) RegisterOption {
	return func(r *Register) {
		r.isPublicEndpoint = b
	}
}
