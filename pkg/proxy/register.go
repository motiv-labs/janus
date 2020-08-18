package proxy

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hellofresh/stats-go/client"
	log "github.com/sirupsen/logrus"
	"go.opencensus.io/plugin/ochttp"

	"github.com/hellofresh/janus/pkg/proxy/balancer"
	"github.com/hellofresh/janus/pkg/proxy/transport"
	"github.com/hellofresh/janus/pkg/router"
)

const (
	methodAll = "ALL"
)

// Register handles the register of proxies into the chosen router.
// It also handles the conversion from a proxy to an http.HandlerFunc
type Register struct {
	router                 router.Router
	idleConnectionsPerHost int
	idleConnTimeout        time.Duration
	idleConnPurgeTicker    *time.Ticker
	flushInterval          time.Duration
	statsClient            client.Client
	matcher                *router.ListenPathMatcher
	isPublicEndpoint       bool
}

// NewRegister creates a new instance of Register
func NewRegister(opts ...RegisterOption) *Register {
	r := Register{
		matcher: router.NewListenPathMatcher(),
	}

	for _, opt := range opts {
		opt(&r)
	}

	return &r
}

// UpdateRouter updates the reference to the router. This is useful to reload the mux
func (p *Register) UpdateRouter(router router.Router) {
	p.router = router
}

// Add register a new route
func (p *Register) Add(definition *RouterDefinition) error {
	log.WithField("balancing_alg", definition.Upstreams.Balancing).Debug("Using a load balancing algorithm")
	balancerInstance, err := balancer.New(definition.Upstreams.Balancing)
	if err != nil {
		log.WithError(err).Error("Could not create a balancer")
		return fmt.Errorf("could not create a balancer: %w", err)
	}

	handler := NewBalancedReverseProxy(definition.Definition, balancerInstance, p.statsClient)
	handler.FlushInterval = p.flushInterval
	handler.Transport = &ochttp.Transport{
		Base: transport.New(
			transport.WithIdleConnTimeout(p.idleConnTimeout),
			transport.WithIdleConnPurgeTicker(p.idleConnPurgeTicker),
			transport.WithInsecureSkipVerify(definition.InsecureSkipVerify),
			transport.WithDialTimeout(time.Duration(definition.ForwardingTimeouts.DialTimeout)),
			transport.WithResponseHeaderTimeout(time.Duration(definition.ForwardingTimeouts.ResponseHeaderTimeout)),
		),
	}

	if p.matcher.Match(definition.ListenPath) {
		p.doRegister(p.matcher.Extract(definition.ListenPath), definition, &ochttp.Handler{Handler: handler, IsPublicEndpoint: p.isPublicEndpoint})
	}

	p.doRegister(definition.ListenPath, definition, &ochttp.Handler{Handler: handler, IsPublicEndpoint: p.isPublicEndpoint})
	return nil
}

func (p *Register) doRegister(listenPath string, def *RouterDefinition, handler http.Handler) {
	log.WithFields(log.Fields{
		"listen_path": listenPath,
	}).Debug("Registering a route")

	if strings.Index(listenPath, "/") != 0 {
		log.WithField("listen_path", listenPath).
			Error("Route listen path must begin with '/'. Skipping invalid route.")
	} else {
		for _, method := range def.Methods {
			if strings.ToUpper(method) == methodAll {
				p.router.Any(listenPath, handler.ServeHTTP, def.middleware...)
			} else {
				p.router.Handle(strings.ToUpper(method), listenPath, handler.ServeHTTP, def.middleware...)
			}
		}
	}
}
