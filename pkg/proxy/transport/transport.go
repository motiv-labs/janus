package transport

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"golang.org/x/net/http2"
)

const (
	// DefaultDialTimeout when connecting to a backend server.
	DefaultDialTimeout = 30 * time.Second

	// DefaultIdleConnsPerHost the default value set for http.Transport.MaxIdleConnsPerHost.
	DefaultIdleConnsPerHost = 64

	// DefaultIdleConnTimeout is the default value for the the maximum amount of time an idle
	// (keep-alive) connection will remain idle before closing itself.
	DefaultIdleConnTimeout = 90 * time.Second
)

type transport struct {
	// Same as net/http.Transport.MaxIdleConnsPerHost, but the default
	// is 64. This value supports scenarios with relatively few remote
	// hosts. When the routing table contains different hosts in the
	// range of hundreds, it is recommended to set this options to a
	// lower value.
	idleConnectionsPerHost int
	insecureSkipVerify     bool
	dialTimeout            time.Duration
	responseHeaderTimeout  time.Duration
	idleConnTimeout        time.Duration
}

// New creates a new instance of Transport with the given params
func New(opts ...Option) *http.Transport {
	t := transport{}

	for _, opt := range opts {
		opt(&t)
	}

	if t.dialTimeout <= 0 {
		t.dialTimeout = DefaultDialTimeout
	}

	if t.idleConnectionsPerHost <= 0 {
		t.idleConnectionsPerHost = DefaultIdleConnsPerHost
	}

	if t.idleConnTimeout == 0 {
		t.idleConnTimeout = DefaultIdleConnTimeout
	}

	tr := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   t.dialTimeout,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       t.idleConnTimeout,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ResponseHeaderTimeout: t.responseHeaderTimeout,
		MaxIdleConnsPerHost:   t.idleConnectionsPerHost,
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: t.insecureSkipVerify},
	}

	http2.ConfigureTransport(tr)

	return tr
}
