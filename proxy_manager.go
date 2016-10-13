package main

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

// ProxyRegister represents a register proxy
type ProxyRegister struct {
	Engine  *gin.Engine
	proxies []Proxy
}

func (p *ProxyRegister) registerMany(proxies []Proxy, breaker *ExtendedCircuitBreakerMeta, beforeHandlers []gin.HandlerFunc, afterHandlers []gin.HandlerFunc) {
	for _, proxy := range proxies {
		p.Register(proxy, breaker, beforeHandlers, afterHandlers)
	}
}

// Register register a new proxy
func (p *ProxyRegister) Register(proxy Proxy, breaker *ExtendedCircuitBreakerMeta, beforeHandlers []gin.HandlerFunc, afterHandlers []gin.HandlerFunc) {
	var handlers []gin.HandlerFunc

	defaultHandler := []gin.HandlerFunc{p.ToHandler(proxy, breaker)}
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

		p.Engine.Any(proxy.ListenPath, handlers...)
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
func (p *ProxyRegister) ToHandler(proxy Proxy, breaker *ExtendedCircuitBreakerMeta) gin.HandlerFunc {
	return func(c *gin.Context) {
		transport := &transport{http.DefaultTransport, breaker, c}
		reverseProxy := NewSingleHostReverseProxy(proxy, transport)
		reverseProxy.ServeHTTP(c.Writer, c.Request)
	}
}
