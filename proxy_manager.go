package janus

import (
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/alexcesaro/statsd.v2"
)

const (
	methodAll = "ALL"
)

// ProxyRegister represents a register proxy
type ProxyRegister struct {
	Router       Router
	proxies      []Proxy
	statsdClient *statsd.Client
}

// RegisterMany registers many proxies at once
func (p *ProxyRegister) RegisterMany(proxies []Proxy, beforeHandlers []HandlerFunc, afterHandlers []HandlerFunc) {
	for _, proxy := range proxies {
		p.Register(proxy, beforeHandlers, afterHandlers)
	}
}

// Register register a new proxy
func (p *ProxyRegister) Register(proxy Proxy, beforeHandlers []HandlerFunc, afterHandlers []HandlerFunc) {
	handler := p.ToHandler(proxy)
	p.registerHandlers(beforeHandlers)

	if false == p.Exists(proxy) {
		log.WithFields(log.Fields{
			"listen_path": proxy.ListenPath,
		}).Info("Registering a proxy")

		for _, method := range proxy.Methods {
			if strings.ToUpper(method) == methodAll {
				p.Router.Any(proxy.ListenPath, handler)
			}

			p.Router.Handle(strings.ToUpper(method), proxy.ListenPath, handler)
		}

		p.proxies = append(p.proxies, proxy)
	}

	p.registerHandlers(afterHandlers)
}

// Exists checks if a proxy is already registered in the manager
func (p *ProxyRegister) Exists(proxy Proxy) bool {
	for _, route := range p.proxies {
		if route.ListenPath == proxy.ListenPath {
			return true
		}
	}

	return false
}

// ToHandler turns a proxy configuration into a handler
func (p *ProxyRegister) ToHandler(proxy Proxy) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		transport := &transport{http.DefaultTransport, p.statsdClient}
		reverseProxy := NewSingleHostReverseProxy(proxy, transport)
		reverseProxy.ServeHTTP(rw, r)
	}
}

func (p *ProxyRegister) registerHandlers(handlers []HandlerFunc) {
	for _, handler := range handlers {
		p.Router.Use(handler)
	}
}
