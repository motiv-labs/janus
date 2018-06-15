package transport

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"strings"
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

func (t transport) hash() string {
	return strings.Join([]string{
		fmt.Sprintf("idleConnectionsPerHost:%v;", t.idleConnectionsPerHost),
		fmt.Sprintf("insecureSkipVerify:%v;", t.insecureSkipVerify),
		fmt.Sprintf("dialTimeout:%v", t.dialTimeout),
		fmt.Sprintf("responseHeaderTimeout:%v", t.responseHeaderTimeout),
		fmt.Sprintf("idleConnTimeout:%v", t.idleConnTimeout),
	}, ";")
}

var registryInstance *registry

func init() {
	registryInstance = newRegistry()
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

	// let's try to get the cached transport from registry, since there is no need to create lots of
	// transports with the same configuration
	hash := t.hash()
	if tr, ok := registryInstance.get(hash); ok {
		return tr
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

	// save newly created transport in registry, to try to reuse it in the future
	registryInstance.put(hash, tr)

	return tr
}
