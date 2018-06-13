package server

import (
	"github.com/hellofresh/janus/pkg/api"
	"github.com/hellofresh/janus/pkg/config"
	"github.com/hellofresh/stats-go/client"
)

// Option represents the available options
type Option func(*Server)

// WithGlobalConfig sets the global configuration
func WithGlobalConfig(globalConfig *config.Specification) Option {
	return func(s *Server) {
		s.globalConfig = globalConfig
	}
}

// WithMetricsClient sets the metric provider
func WithMetricsClient(client client.Client) Option {
	return func(s *Server) {
		s.statsClient = client
	}
}

// WithProvider sets the configuration provider
func WithProvider(provider api.Repository) Option {
	return func(s *Server) {
		s.provider = provider
	}
}

// WithProfiler enables or disables profiler
func WithProfiler(enabled, public bool) Option {
	return func(s *Server) {
		s.profilingEnabled = enabled
		s.profilingPublic = public
	}
}
