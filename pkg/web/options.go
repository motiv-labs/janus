package web

import (
	"github.com/hellofresh/janus/pkg/api"
	"github.com/hellofresh/janus/pkg/config"
)

// Option represents the available options
type Option func(*Server)

// WithConfigurations sets the current configurations in memory
func WithConfigurations(cfgs *api.Configuration) Option {
	return func(s *Server) {
		s.apiHandler.Cfgs = cfgs
	}
}

// WithPort sets the server port
func WithPort(port int) Option {
	return func(s *Server) {
		s.Port = port
	}
}

// WithCredentials sets the credentials for the server
func WithCredentials(cred config.Credentials) Option {
	return func(s *Server) {
		s.Credentials = cred
	}
}

// WithTLS sets the TLS configs for the server
func WithTLS(tls config.TLS) Option {
	return func(s *Server) {
		s.TLS = tls
	}
}
