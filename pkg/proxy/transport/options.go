package transport

import (
	"time"
)

// Option represents the transport options
type Option func(*transport)

// WithCloseIdleConnsPeriod sets the time period of how often the idle connections are
// forcibly closed
func WithCloseIdleConnsPeriod(d time.Duration) Option {
	return func(t *transport) {
		t.closeIdleConnsPeriod = d
	}
}

// WithInsecureSkipVerify sets tls config insecure skip verify
func WithInsecureSkipVerify(value bool) Option {
	return func(t *transport) {
		t.insecureSkipVerify = value
	}
}

// WithDialTimeout sets the dial context timeout
func WithDialTimeout(d time.Duration) Option {
	return func(t *transport) {
		t.dialTimeout = d
	}
}

// WithResponseHeaderTimeout sets the response header timeout
func WithResponseHeaderTimeout(d time.Duration) Option {
	return func(t *transport) {
		t.responseHeaderTimeout = d
	}
}
