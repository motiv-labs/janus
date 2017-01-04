package proxy

import (
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/hellofresh/janus/pkg/router"
)

const (
	methodAll = "ALL"
)

type Register interface {
	Exists(route *Route) bool
	Get(listenPath string) *Route
	Remove(listenPath string) error
	AddMany(routes []*Route) error
	Add(route *Route) error
}

// Register handles the register of proxies into the choosen router.
// It also handles the conversion from a proxy to an http.HandlerFunc
type InMemoryRegister struct {
	router    router.Router
	transport http.RoundTripper
	proxies   map[string]*Route
}

// NewInMemoryRegister creates a new instance of Register
func NewInMemoryRegister(router router.Router, transport http.RoundTripper) *InMemoryRegister {
	return &InMemoryRegister{router: router, transport: transport, proxies: make(map[string]*Route)}
}

// Exists checks if a proxy is already registered in the manager
func (p *InMemoryRegister) Exists(route *Route) bool {
	_, exists := p.proxies[route.proxy.ListenPath]
	return exists
}

func (p *InMemoryRegister) Get(listenPath string) *Route {
	return p.proxies[listenPath]
}

func (p *InMemoryRegister) Remove(listenPath string) error {
	delete(p.proxies, listenPath)
	return nil
}

// AddMany registers many proxies at once
func (p *InMemoryRegister) AddMany(routes []*Route) error {
	for _, r := range routes {
		err := p.Add(r)
		if nil != err {
			return err
		}
	}

	return nil
}

// Register register a new proxy
func (p *InMemoryRegister) Add(route *Route) error {
	if p.Exists(route) {
		return p.replace(route)
	} else {
		return p.add(route)
	}
}

func (p *InMemoryRegister) add(route *Route) error {
	if false == p.Exists(route) {
		proxy := route.proxy

		handler := p.toHandler(proxy)
		matcher := router.NewListenPathMatcher()
		if matcher.Match(proxy.ListenPath) {
			p.doRegister(matcher.Extract(proxy.ListenPath), handler, proxy.Methods, route.handlers)
		}

		p.doRegister(proxy.ListenPath, handler, proxy.Methods, route.handlers)
		p.proxies[proxy.ListenPath] = route
	}

	return nil
}

func (p *InMemoryRegister) replace(r *Route) error {
	log.WithFields(log.Fields{
		"listen_path": r.proxy.ListenPath,
		"target_url":  r.proxy.TargetURL,
	}).Debug("Replacing a route")

	currentRoute := p.Get(r.proxy.ListenPath)
	*currentRoute.proxy = *r.proxy
	currentRoute.handlers = r.handlers

	return nil
}

func (p *InMemoryRegister) doRegister(listenPath string, handler http.HandlerFunc, methods []string, handlers []router.Constructor) {
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

// ToHandler turns a proxy configuration into a handler
func (p *InMemoryRegister) toHandler(proxy *Proxy) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		reverseProxy := NewSingleHostReverseProxy(proxy, p.transport)
		reverseProxy.ServeHTTP(rw, r)
	}
}
