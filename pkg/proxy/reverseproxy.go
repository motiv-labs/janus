package proxy

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/hellofresh/janus/pkg/router"
	stats "github.com/hellofresh/stats-go"
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
	// StatsClient defines the stats client for tracing
	StatsClient stats.Client

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

// Proxy instances implement Janus proxying functionality. For
// initializing, see the WithParams the constructor and Params.
type Proxy struct {
	statsClient   stats.Client
	quit          chan struct{}
	flushInterval time.Duration
}

// WithParams returns a new ReverseProxy that routes
// URLs to the scheme, host, and base path provided in target. If the
// target's path is "/base" and the incoming request was for "/dir",
// the target request will be for /base/dir.
// NewSingleHostReverseProxy does not rewrite the Host header.
// To rewrite Host headers, use ReverseProxy directly with a custom
// Director policy.
func WithParams(o Params) *Proxy {
	if o.IdleConnectionsPerHost <= 0 {
		o.IdleConnectionsPerHost = DefaultIdleConnsPerHost
	}

	if o.CloseIdleConnsPeriod == 0 {
		o.CloseIdleConnsPeriod = DefaultCloseIdleConnsPeriod
	}

	tr := http.DefaultTransport.(*http.Transport)
	tr.Proxy = http.ProxyFromEnvironment
	tr.DialContext = (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}).DialContext
	tr.MaxIdleConns = 100
	tr.IdleConnTimeout = 90 * time.Second
	tr.TLSHandshakeTimeout = 10 * time.Second
	tr.ExpectContinueTimeout = 1 * time.Second
	tr.MaxIdleConnsPerHost = o.IdleConnectionsPerHost

	quit := make(chan struct{})
	if o.CloseIdleConnsPeriod > 0 {
		go func() {
			for {
				select {
				case <-time.After(o.CloseIdleConnsPeriod):
					tr.CloseIdleConnections()
				case <-quit:
					return
				}
			}
		}()
	}

	if o.InsecureSkipVerify {
		tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	return &Proxy{
		statsClient:   o.StatsClient,
		quit:          quit,
		flushInterval: o.FlushInterval,
	}
}

// Reverse given a target and chains of inbound/outbound plugins, we make a ReverseProxy
func (p *Proxy) Reverse(proxyDefinition *Definition, inbound InChain, outbound OutChain) *httputil.ReverseProxy {
	director := func(req *http.Request) {
		target, _ := url.Parse(proxyDefinition.UpstreamURL)
		targetQuery := target.RawQuery

		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		path := target.Path

		if proxyDefinition.AppendPath {
			log.Debug("Appending listen path to the target url")
			path = singleJoiningSlash(target.Path, req.URL.Path)
		}

		if proxyDefinition.StripPath {
			path = singleJoiningSlash(target.Path, req.URL.Path)
			matcher := router.NewListenPathMatcher()
			listenPath := matcher.Extract(proxyDefinition.ListenPath)

			log.Debugf("Stripping listen path: %s", listenPath)
			path = strings.Replace(path, listenPath, "", 1)
			if !strings.HasSuffix(target.Path, "/") && strings.HasSuffix(path, "/") {
				path = path[:len(path)-1]
			}
		}

		log.Debugf("Upstream Path is: %s", path)
		req.URL.Path = path

		// This is very important to avoid problems with ssl verification for the HOST header
		if !proxyDefinition.PreserveHost {
			log.Debug("Preserving the host header")
			req.Host = target.Host
		}

		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
	}

	transport := &http.Transport{
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
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: proxyDefinition.InsecureSkipVerify,
		},
	}

	return &httputil.ReverseProxy{
		Director:      director,
		Transport:     &Shackles{p.statsClient, inbound, outbound, transport},
		FlushInterval: p.flushInterval,
	}
}

// Close causes the proxy to stop closing idle
// connections and, currently, has no other effect.
// It's primary purpose is to support testing.
func (p *Proxy) Close() error {
	close(p.quit)
	return nil
}

func cleanSlashes(a string) string {
	endSlash := strings.HasSuffix(a, "//")
	startSlash := strings.HasPrefix(a, "//")

	if startSlash {
		a = "/" + strings.TrimPrefix(a, "//")
	}

	if endSlash {
		a = strings.TrimSuffix(a, "//") + "/"
	}

	return a
}

func singleJoiningSlash(a, b string) string {
	a = cleanSlashes(a)
	b = cleanSlashes(b)

	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")

	switch {
	case aslash && bslash:
		log.Debug(a + b)
		return a + b[1:]
	case !aslash && !bslash:
		if len(b) > 0 {
			log.Debug(a + b)
			return a + "/" + b
		}

		log.Debug(a + b)
		return a
	}

	log.Debug(a + b)
	return a + b
}
