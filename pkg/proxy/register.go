package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
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
	params Params
}

// NewRegister creates a new instance of Register
func NewRegister(router router.Router, params Params) *Register {
	return &Register{router, params}
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
	definition := route.Proxy

	p.params.Outbound = route.Outbound
	p.params.InsecureSkipVerify = definition.InsecureSkipVerify
	handler := &httputil.ReverseProxy{
		Director:  p.createDirector(definition),
		Transport: NewTransportWithParams(p.params),
	}

	matcher := router.NewListenPathMatcher()
	if matcher.Match(definition.ListenPath) {
		p.doRegister(matcher.Extract(definition.ListenPath), handler.ServeHTTP, definition.Methods, route.Inbound)
	}

	p.doRegister(definition.ListenPath, handler.ServeHTTP, definition.Methods, route.Inbound)
	return nil
}

func (p *Register) createDirector(proxyDefinition *Definition) func(req *http.Request) {
	return func(req *http.Request) {
		target, _ := url.Parse(proxyDefinition.UpstreamURL)
		targetQuery := target.RawQuery

		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		path := target.Path

		if proxyDefinition.AppendPath {
			log.Debug("Appending listen path to the target url")
			path = singleJoiningSlash(target.Path, req.URL.Path)
		}

		if proxyDefinition.StripPath {
			path = singleJoiningSlash(target.Path, req.URL.Path)
			matcher := router.NewListenPathMatcher()
			listenPath := matcher.Extract(proxyDefinition.ListenPath)

			log.WithField("listen_path", listenPath).Debug("Stripping listen path")
			path = strings.Replace(path, listenPath, "", 1)
			if !strings.HasSuffix(target.Path, "/") && strings.HasSuffix(path, "/") {
				path = path[:len(path)-1]
			}
		}

		log.Debugf("Upstream Path is: %s", path)
		req.URL.Path = path

		// This is very important to avoid problems with ssl verification for the HOST header
		if !proxyDefinition.PreserveHost {
			log.Debug("Preserving the host header")
			req.Host = target.Host
		}

		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
	}
}

func (p *Register) doRegister(listenPath string, handler http.HandlerFunc, methods []string, handlers InChain) {
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

func cleanSlashes(a string) string {
	endSlash := strings.HasSuffix(a, "//")
	startSlash := strings.HasPrefix(a, "//")

	if startSlash {
		a = "/" + strings.TrimPrefix(a, "//")
	}

	if endSlash {
		a = strings.TrimSuffix(a, "//") + "/"
	}

	return a
}

func singleJoiningSlash(a, b string) string {
	a = cleanSlashes(a)
	b = cleanSlashes(b)

	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")

	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		if len(b) > 0 {
			return a + "/" + b
		}
		return a
	}
	return a + b
}
