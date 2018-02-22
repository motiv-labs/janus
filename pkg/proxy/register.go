package proxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/hellofresh/janus/pkg/router"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const (
	methodAll = "ALL"
)

// Register handles the register of proxies into the chosen router.
// It also handles the conversion from a proxy to an http.HandlerFunc
type Register struct {
	Router router.Router
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

	var balancer Balancer
	if definition.IsBalancerDefined() {
		log.WithField("balancing_alg", definition.Upstreams.Balancing).Debug("Using a load balancing algorithm")

		var err error
		balancer, err = NewBalancer(definition.Upstreams.Balancing)
		if err != nil {
			msg := "Could not create a balancer"
			log.WithError(err).Error(msg)
			return errors.Wrap(err, msg)
		}
	}

	p.params.Breaker.ID = definition.ListenPath
	p.params.Outbound = route.Outbound
	p.params.InsecureSkipVerify = definition.InsecureSkipVerify
	handler := &httputil.ReverseProxy{
		Director:  p.createDirector(definition, balancer),
		Transport: NewTransportWithParams(p.params),
	}

	matcher := router.NewListenPathMatcher()
	if matcher.Match(definition.ListenPath) {
		p.doRegister(matcher.Extract(definition.ListenPath), handler.ServeHTTP, definition.Methods, route.Inbound)
	}

	p.doRegister(definition.ListenPath, handler.ServeHTTP, definition.Methods, route.Inbound)
	return nil
}

func (p *Register) createDirector(proxyDefinition *Definition, balancer Balancer) func(req *http.Request) {
	paramNameExtractor := router.NewListenPathParamNameExtractor()
	matcher := router.NewListenPathMatcher()

	return func(req *http.Request) {
		var upstreamURL string
		if proxyDefinition.IsBalancerDefined() && balancer != nil {
			upstream, err := balancer.Elect(proxyDefinition.Upstreams.Targets)
			if err != nil {
				log.WithError(err).Error("Could not elect one upstream")
				return
			}
			log.WithField("target", upstream.Target).Debug("Target upstream elected")
			upstreamURL = upstream.Target
		} else {
			log.WithField("upstream_url", proxyDefinition.UpstreamURL).
				Warn("The upstream URL is deprecated. Use Upstreams instead")
			upstreamURL = proxyDefinition.UpstreamURL
		}

		target, err := url.Parse(upstreamURL)
		if err != nil {
			log.WithError(err).WithField("upstream_url", upstreamURL).Error("Could not parse the target URL")
			return
		}

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
			listenPath := matcher.Extract(proxyDefinition.ListenPath)

			log.WithField("listen_path", listenPath).Debug("Stripping listen path")
			path = strings.Replace(path, listenPath, "", 1)
			if !strings.HasSuffix(target.Path, "/") && strings.HasSuffix(path, "/") {
				path = path[:len(path)-1]
			}
		}

		paramNames := paramNameExtractor.Extract(path)
		parametrizedPath, err := p.applyParameters(req, path, paramNames)
		if err != nil {
			log.WithError(err).Warn("Unable to extract param from request")
		} else {
			path = parametrizedPath
		}

		log.WithField("path", path).Debug("Upstream Path")
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

	if strings.Index(listenPath, "/") != 0 {
		log.WithField("listen_path", listenPath).
			Error("Route listen path must begin with '/'. Skipping invalid route.")
	} else {
		for _, method := range methods {
			if strings.ToUpper(method) == methodAll {
				p.Router.Any(listenPath, handler, handlers...)
			} else {
				p.Router.Handle(strings.ToUpper(method), listenPath, handler, handlers...)
			}
		}
	}
}

func (p *Register) applyParameters(req *http.Request, path string, paramNames []string) (string, error) {
	for _, paramName := range paramNames {
		paramValue := router.URLParam(req, paramName)

		if len(paramValue) == 0 {
			return "", errors.Errorf("unable to extract {%s} from request", paramName)
		}

		path = strings.Replace(
			path,
			fmt.Sprintf("{%s}", paramName),
			paramValue,
			-1,
		)
	}

	return path, nil
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

	aSlash := strings.HasSuffix(a, "/")
	bSlash := strings.HasPrefix(b, "/")

	switch {
	case aSlash && bSlash:
		return a + b[1:]
	case !aSlash && !bSlash:
		if len(b) > 0 {
			return a + "/" + b
		}
		return a
	}
	return a + b
}
