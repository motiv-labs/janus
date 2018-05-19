package proxy

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"
)

const (
	// DefaultIdleConnsPerHost the default value set for http.Transport.MaxIdleConnsPerHost.
	DefaultIdleConnsPerHost = 64

	// DefaultCloseIdleConnsPeriod the default period at which the idle connections are forcibly
	// closed.
	DefaultCloseIdleConnsPeriod = 20 * time.Second
)

// Params initialization options.
type Params struct {
	// When set, the proxy will skip the TLS verification on outgoing requests.
	InsecureSkipVerify bool

	// Same as net/http.Transport.MaxIdleConnsPerHost, but the default
	// is 64. This value supports scenarios with relatively few remote
	// hosts. When the routing table contains different hosts in the
	// range of hundreds, it is recommended to set this options to a
	// lower value.
	IdleConnectionsPerHost int

	// Defines the time period of how often the idle connections are
	// forcibly closed. The default is 12 seconds. When set to less than
	// 0, the proxy doesn't force closing the idle connections.
	CloseIdleConnsPeriod time.Duration

	// The Flush interval for copying upgraded connections
	FlushInterval time.Duration
}

// NewTransportWithParams creates a new instance of Transport with the given params
func NewTransportWithParams(o Params) *http.Transport {
	if o.IdleConnectionsPerHost <= 0 {
		o.IdleConnectionsPerHost = DefaultIdleConnsPerHost
	}

	if o.CloseIdleConnsPeriod == 0 {
		o.CloseIdleConnsPeriod = DefaultCloseIdleConnsPeriod
	}

	tr := http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		MaxIdleConnsPerHost:   o.IdleConnectionsPerHost,
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: o.InsecureSkipVerify},
	}

	if o.CloseIdleConnsPeriod > 0 {
		go func() {
			for range time.After(o.CloseIdleConnsPeriod) {
				tr.CloseIdleConnections()
			}
		}()
	}

	return &tr
}
