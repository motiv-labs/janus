package proxy

import (
	"net/http"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/hellofresh/janus/log"
	"github.com/hellofresh/janus/router"
)

const (
	methodAll = "ALL"
)

// RegisterChan holds two channels. Many is a channel of many routes
// this is more for a bulk operation, when you create many routes at once
// use this channel.
// One is for a single route, if you create one route and want to register it
// use this channel
type RegisterChan struct {
	Many chan []*Route
	One  chan *Route
}

// NewRegisterChan creates a new instance of RegisterChan
func NewRegisterChan(router router.Router, transport http.RoundTripper) *RegisterChan {
	register := NewRegister(router, transport)

	registerChan := &RegisterChan{}
	registerChan.Many = make(chan []*Route)
	registerChan.One = make(chan *Route)

	go registerChan.listenForRouteChanges(register)
	return registerChan
}

func (rc *RegisterChan) listenForRouteChanges(register *Register) {
	for {
		select {
		case routes := <-rc.Many:
			register.RegisterMany(routes)

		case route := <-rc.One:
			register.Register(route)
		}
	}
}

// Register handles the register of proxies into the choosen router.
// It also handles the conversion from a proxy to an http.HandlerFunc
type Register struct {
	router    router.Router
	transport http.RoundTripper
	proxies   []Proxy
}

// NewRegister creates a new instance of Register
func NewRegister(router router.Router, transport http.RoundTripper) *Register {
	return &Register{router: router, transport: transport}
}

// RegisterMany registers many proxies at once
func (p *Register) RegisterMany(routes []*Route) {
	for _, route := range routes {
		p.Register(route)
	}
}

// Register register a new proxy
func (p *Register) Register(route *Route) {
	proxy := route.proxy

	if false == p.Exists(proxy) {
		handler := p.ToHandler(proxy)
		matcher := router.NewListenPathMatcher()
		if matcher.Match(proxy.ListenPath) {
			p.doRegister(matcher.Extract(proxy.ListenPath), handler, proxy.Methods, route.handlers)
		}

		p.doRegister(proxy.ListenPath, handler, proxy.Methods, route.handlers)
		p.proxies = append(p.proxies, proxy)
	}

}

func (p *Register) doRegister(
	listenPath string,
	handler http.HandlerFunc,
	methods []string,
	handlers []router.Constructor,
) {
	log.WithFields(logrus.Fields{
		"listen_path": listenPath,
	}).Debug("Registering a proxy")

	for _, method := range methods {
		if strings.ToUpper(method) == methodAll {
			p.router.Any(listenPath, handler, handlers...)
		} else {
			p.router.Handle(strings.ToUpper(method), listenPath, handler, handlers...)
		}
	}
}

// Exists checks if a proxy is already registered in the manager
func (p *Register) Exists(proxy Proxy) bool {
	for _, route := range p.proxies {
		if route.ListenPath == proxy.ListenPath {
			return true
		}
	}

	return false
}

// ToHandler turns a proxy configuration into a handler
func (p *Register) ToHandler(proxy Proxy) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		reverseProxy := NewSingleHostReverseProxy(proxy, p.transport)
		reverseProxy.ServeHTTP(rw, r)
	}
}
