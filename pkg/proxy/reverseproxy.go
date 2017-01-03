package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/hellofresh/janus/pkg/request"
	"github.com/hellofresh/janus/pkg/router"
)

var (
	// ContextKeyBody defines the db context key
	ContextKeyBody = request.ContextKey("body")
)

// NewSingleHostReverseProxy returns a new ReverseProxy that routes
// URLs to the scheme, host, and base path provided in target. If the
// target's path is "/base" and the incoming request was for "/dir",
// the target request will be for /base/dir.
// NewSingleHostReverseProxy does not rewrite the Host header.
// To rewrite Host headers, use ReverseProxy directly with a custom
// Director policy.
func NewSingleHostReverseProxy(proxy Proxy, transport http.RoundTripper) *httputil.ReverseProxy {
	target, _ := url.Parse(proxy.TargetURL)
	targetQuery := target.RawQuery

	director := func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		path := target.Path

		if proxy.AppendListenPath {
			log.Debug("Appending listen path to the target url")
			path = singleJoiningSlash(target.Path, req.URL.Path)
		}

		if proxy.StripListenPath {
			path = singleJoiningSlash(target.Path, req.URL.Path)
			matcher := router.NewListenPathMatcher()
			listenPath := matcher.Extract(proxy.ListenPath)

			log.Debugf("Stripping listen path: %s", listenPath)
			path = strings.Replace(path, listenPath, "", 1)
			if !strings.HasSuffix(target.Path, "/") && strings.HasSuffix(path, "/") {
				path = path[:len(path)-1]
			}
		}

		log.Debugf("Upstream Path is: %s", path)
		req.URL.Path = path

		// This is very important to avoid problems with ssl verification for the HOST header
		if !proxy.PreserveHostHeader {
			log.Debug("Preserving the host header")
			req.Host = target.Host
		}

		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
	}

	return &httputil.ReverseProxy{Director: director, Transport: transport}
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
		log.Debug(a + b)
		return a + b[1:]
	case !aslash && !bslash:
		if len(b) > 0 {
			log.Debug(a + b)
			return a + "/" + b
		}

		log.Debug(a + b)
		return a
	}

	log.Debug(a + b)
	return a + b
}
