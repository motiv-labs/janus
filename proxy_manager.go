package main

import (
	"github.com/kataras/iris"
	"net/url"
	"net/http/httputil"
	"strings"
	"net/http"
	log "github.com/Sirupsen/logrus"
)

type transport struct {
	http.RoundTripper
	breaker ExtendedCircuitBreakerMeta
}

func (t *transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {

	if t.breaker.CB.Ready() {
		log.Debug("ON REQUEST: Breaker status: ", t.breaker.CB.Ready())
		resp, err = t.RoundTripper.RoundTrip(req)

		if err != nil {
			t.breaker.CB.Fail()
		} else if resp.StatusCode == 500 {
			t.breaker.CB.Fail()
		} else {
			t.breaker.CB.Success()
		}
	}
	
	return resp, nil
}

var _ http.RoundTripper = &transport{}

type ProxyRegister struct{}

func NewProxyRegister() *ProxyRegister {
	return &ProxyRegister{}
}

func (p *ProxyRegister) RegisterMany(proxies []Proxy, breaker ExtendedCircuitBreakerMeta) {
	for _, proxy := range proxies {
		p.Register(proxy, breaker)
	}
}

func (p *ProxyRegister) Register(proxy Proxy, breaker ExtendedCircuitBreakerMeta) {
	handler := p.createHandler(proxy, breaker)

	iris.Handle("", proxy.ListenPath, iris.ToHandler(handler))
}

func (p *ProxyRegister) createHandler(proxy Proxy, breaker ExtendedCircuitBreakerMeta) *httputil.ReverseProxy {
	target, _ := url.Parse(proxy.TargetURL)

	director := func(req *http.Request) {
		log.Debug("Started proxy")
		path := target.Path
		targetQuery := target.RawQuery

		if proxy.StripListenPath {
			log.Debug("Stripping: ", proxy.ListenPath)
			listenPath := strings.Replace(proxy.ListenPath, "/*randomName", "", -1)

			path = singleJoiningSlash(target.Path, req.URL.Path)
			path = strings.Replace(path, listenPath, "", -1)

			log.Debug("Upstream Path is: ", path)
		}

		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = path

		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}

		log.Debug("Done proxy")
	}

	return &httputil.ReverseProxy{Director: director, Transport: &transport{http.DefaultTransport, breaker}}
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}
