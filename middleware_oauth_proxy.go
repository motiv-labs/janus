package main

import (
	"github.com/valyala/fasthttp"
	"github.com/kataras/iris"
)

// a silly example
type OauthProxy struct {
	*Middleware
	proxyRegister *ProxyRegister
}

//Important staff, iris middleware must implement the iris.Handler interface which is:
func (m OauthProxy) ProcessRequest(req fasthttp.Request, resp fasthttp.Response, c *iris.Context) (error, int) {
	if m.Spec.UseOauth2 {
		return nil, fasthttp.StatusOK
	}

	var proxies []Proxy

	oauthMeta := m.Spec.Oauth2Meta

	//oauth proxy
	proxies = append(proxies, oauthMeta.OauthEndpoints.Authorize)
	proxies = append(proxies, oauthMeta.OauthEndpoints.Token)
	proxies = append(proxies, oauthMeta.OauthEndpoints.Info)
	proxies = append(proxies, oauthMeta.OauthClientEndpoints.Create)
	proxies = append(proxies, oauthMeta.OauthClientEndpoints.Remove)

	cb := NewCircuitBreaker(m.Spec)
	m.proxyRegister.registerMany(proxies, cb)

	return nil, fasthttp.StatusOK
}
