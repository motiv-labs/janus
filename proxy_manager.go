package main

import (
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"gopkg.in/alexcesaro/statsd.v2"
)

const (
	methodAll = "ALL"
)

// ProxyRegister represents a register proxy
type ProxyRegister struct {
	Engine       *gin.Engine
	proxies      []Proxy
	statsdClient *statsd.Client
}

// RegisterMany registers many proxies at once
func (p *ProxyRegister) RegisterMany(proxies []Proxy, beforeHandlers []gin.HandlerFunc, afterHandlers []gin.HandlerFunc) {
	for _, proxy := range proxies {
		p.Register(proxy, beforeHandlers, afterHandlers)
	}
}

// Register register a new proxy
func (p *ProxyRegister) Register(proxy Proxy, beforeHandlers []gin.HandlerFunc, afterHandlers []gin.HandlerFunc) {
	var handlers []gin.HandlerFunc

	defaultHandler := []gin.HandlerFunc{p.ToHandler(proxy)}
	handlers = append(defaultHandler, handlers...)

	if len(beforeHandlers) > 0 {
		handlers = append(beforeHandlers, handlers...)
	}

	if len(afterHandlers) > 0 {
		handlers = append(handlers, afterHandlers...)
	}

	if false == p.Exists(proxy) {
		log.WithFields(log.Fields{
			"listen_path": proxy.ListenPath,
		}).Info("Registering a proxy")

		for _, method := range proxy.Methods {
			if strings.ToUpper(method) == methodAll {
				p.Engine.Any(proxy.ListenPath, handlers...)
			}

			p.Engine.Handle(strings.ToUpper(method), proxy.ListenPath, handlers...)
		}

		p.proxies = append(p.proxies, proxy)
	}
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
func (p *ProxyRegister) ToHandler(proxy Proxy) gin.HandlerFunc {
	return func(c *gin.Context) {
		transport := &transport{http.DefaultTransport, c, p.statsdClient}
		reverseProxy := NewSingleHostReverseProxy(proxy, transport)
		reverseProxy.ServeHTTP(c.Writer, c.Request)
	}
}
