package proxy

import (
	"net/http"
	"strings"

	"github.com/hellofresh/janus/pkg/router"
	log "github.com/sirupsen/logrus"
)

const (
	methodAll = "ALL"
)

// Register handles the register of proxies into the chosen router.
// It also handles the conversion from a proxy to an http.HandlerFunc
type Register struct {
	router router.Router
	proxy  *Proxy
}

// NewRegister creates a new instance of Register
func NewRegister(router router.Router, proxy *Proxy) *Register {
	return &Register{router: router, proxy: proxy}
}

// AddMany registers many proxies at once
func (p *Register) AddMany(routes []*Route) error {
	for _, r := range routes {
		err := p.Add(r)
		if nil != err {
			return err
		}
	}

	return nil
}

// Add register a new route
func (p *Register) Add(route *Route) error {
	return p.AddWithInOut(route, InChain{}, OutChain{})
}

// AddWithInOut register a new route with inbound/outbounds plugins
func (p *Register) AddWithInOut(route *Route, inbound InChain, outbound OutChain) error {
	definition := route.proxy

	handler := p.proxy.Reverse(definition, inbound, outbound).ServeHTTP
	matcher := router.NewListenPathMatcher()
	if matcher.Match(definition.ListenPath) {
		p.doRegister(matcher.Extract(definition.ListenPath), handler, definition.Methods, route.handlers)
	}

	p.doRegister(definition.ListenPath, handler, definition.Methods, route.handlers)
	return nil
}

func (p *Register) doRegister(listenPath string, handler http.HandlerFunc, methods []string, handlers []router.Constructor) {
	log.WithFields(log.Fields{
		"listen_path": listenPath,
	}).Debug("Registering a route")

	for _, method := range methods {
		if strings.ToUpper(method) == methodAll {
			p.router.Any(listenPath, handler, handlers...)
		} else {
			p.router.Handle(strings.ToUpper(method), listenPath, handler, handlers...)
		}
	}
}
