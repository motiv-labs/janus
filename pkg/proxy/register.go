package proxy

import (
	"net/http"
	"strings"

	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/hellofresh/janus/pkg/router"
	"github.com/hellofresh/janus/pkg/store"
)

const (
	methodAll = "ALL"
)

// Register handles the register of proxies into the chosen router.
// It also handles the conversion from a proxy to an http.HandlerFunc
type Register struct {
	router router.Router
	proxy  *Proxy
	store  store.Store
}

// NewRegister creates a new instance of Register
func NewRegister(router router.Router, proxy *Proxy, store store.Store) *Register {
	return &Register{router: router, proxy: proxy, store: store}
}

// Exists checks if a proxy is already registered in the manager
func (p *Register) Exists(route *Route) bool {
	exists, _ := p.store.Exists(route.proxy.ListenPath)
	return exists
}

// Get
func (p *Register) Get(listenPath string) *Route {
	var route Route
	rawRoute, err := p.store.Get(listenPath)
	if nil != err {
		log.Warn(err.Error())
	}

	json.Unmarshal([]byte(rawRoute), &route)
	if nil != err {
		log.Warn(err.Error())
	}

	return &route
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
	if p.Exists(route) {
		return p.replace(route)
	}

	return p.add(route)
}

func (p *Register) add(route *Route) error {
	definition := route.proxy

	handler := p.proxy.Reverse(definition).ServeHTTP
	matcher := router.NewListenPathMatcher()
	if matcher.Match(definition.ListenPath) {
		p.doRegister(matcher.Extract(definition.ListenPath), handler, definition.Methods, route.handlers)
	}

	p.doRegister(definition.ListenPath, handler, definition.Methods, route.handlers)
	jsonRoute, err := json.Marshal(route)

	if err != nil {
		return err
	}

	p.store.Set(definition.ListenPath, string(jsonRoute), 0)

	return nil
}

func (p *Register) replace(r *Route) error {
	log.WithFields(log.Fields{
		"listen_path": r.proxy.ListenPath,
		"target_url":  r.proxy.TargetURL,
	}).Debug("Replacing a route")

	currentRoute := p.Get(r.proxy.ListenPath)
	*currentRoute.proxy = *r.proxy
	currentRoute.handlers = r.handlers

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
