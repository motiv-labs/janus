package proxy

import (
	"net/http"
	"strings"
	"time"

	"github.com/hellofresh/janus/pkg/proxy/balancer"
	"github.com/hellofresh/janus/pkg/proxy/transport"
	"github.com/hellofresh/janus/pkg/router"
	"github.com/hellofresh/stats-go/client"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const (
	methodAll = "ALL"
)

// Register handles the register of proxies into the chosen router.
// It also handles the conversion from a proxy to an http.HandlerFunc
type Register struct {
	router                 router.Router
	idleConnectionsPerHost int
	closeIdleConnsPeriod   time.Duration
	flushInterval          time.Duration
	statsClient            client.Client
}

// NewRegister creates a new instance of Register
func NewRegister(opts ...RegisterOption) *Register {
	r := Register{}

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
func (p *Register) Add(definition *Definition) error {
	log.WithField("balancing_alg", definition.Upstreams.Balancing).Debug("Using a load balancing algorithm")
	balancerInstance, err := balancer.New(definition.Upstreams.Balancing)
	if err != nil {
		msg := "Could not create a balancer"
		log.WithError(err).Error(msg)
		return errors.Wrap(err, msg)
	}

	handler := NewBalancedReverseProxy(definition, balancerInstance, p.statsClient)
	handler.FlushInterval = p.flushInterval
	handler.Transport = transport.New(
		transport.WithCloseIdleConnsPeriod(p.closeIdleConnsPeriod),
		transport.WithInsecureSkipVerify(definition.InsecureSkipVerify),
		transport.WithDialTimeout(time.Duration(definition.ForwardingTimeouts.DialTimeout)),
		transport.WithResponseHeaderTimeout(time.Duration(definition.ForwardingTimeouts.ResponseHeaderTimeout)),
	)

	matcher := router.NewListenPathMatcher()
	if matcher.Match(definition.ListenPath) {
		p.doRegister(matcher.Extract(definition.ListenPath), definition, handler.ServeHTTP)
	}

	p.doRegister(definition.ListenPath, definition, handler.ServeHTTP)
	return nil
}

func (p *Register) doRegister(listenPath string, def *Definition, handler http.HandlerFunc) {
	log.WithFields(log.Fields{
		"listen_path": listenPath,
	}).Debug("Registering a route")

	if strings.Index(listenPath, "/") != 0 {
		log.WithField("listen_path", listenPath).
			Error("Route listen path must begin with '/'. Skipping invalid route.")
	} else {
		for _, method := range def.Methods {
			if strings.ToUpper(method) == methodAll {
				p.router.Any(listenPath, handler, def.middleware...)
			} else {
				p.router.Handle(strings.ToUpper(method), listenPath, handler, def.middleware...)
			}
		}
	}
}
