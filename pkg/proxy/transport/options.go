package transport

import (
	"time"
)

// Option represents the transport options
type Option func(*transport)

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

// WithIdleConnTimeout sets the maximum amount of time an idle
// (keep-alive) connection will remain idle before closing
// itself.
func WithIdleConnTimeout(d time.Duration) Option {
	return func(t *transport) {
		t.idleConnTimeout = d
	}
}

// WithIdleConnPurgeTicker purges idle connections on every interval if set
// this is done to prevent permanent keep-alive on connections with high ops
func WithIdleConnPurgeTicker(ticker *time.Ticker) Option {
	return func(t *transport) {
		t.idleConnPurgeTicker = ticker
	}
}
