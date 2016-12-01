package janus

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	statsd "gopkg.in/alexcesaro/statsd.v2"
)

type transport struct {
	http.RoundTripper
	context      *gin.Context
	statsdClient *statsd.Client
}

func (t *transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	timing := t.statsdClient.NewTiming()
	resp, err = t.RoundTripper.RoundTrip(req)
	timing.Send(getStatsdMetricName(req))

	if resp.StatusCode >= 400 {
		t.statsdClient.Increment("error_request")
	} else if resp.StatusCode < 300 && resp.Body != nil {
		t.statsdClient.Increment("success_request")

		//This is useful for the middlewares
		var bodyBytes []byte

		defer resp.Body.Close()
		bodyBytes, _ = ioutil.ReadAll(resp.Body)

		// Restore the io.ReadCloser to its original state
		resp.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		// Use the content
		log.WithFields(log.Fields{
			"req":  req,
			"resp": resp,
		}).Info("Setting body")

		t.context.Set("body", bodyBytes)
	} else {
		t.statsdClient.Increment("success_request")
	}

	return resp, err
}

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

		if proxy.StripListenPath {
			path = singleJoiningSlash(target.Path, req.URL.Path)

			listenPath := extractListenPath(proxy.ListenPath)
			log.Debug("Stripping: ", listenPath)
			path = strings.Replace(path, listenPath, "", 1)

			log.Debug("Upstream Path is: ", path)

			if !strings.HasSuffix(target.Path, "/") && strings.HasSuffix(path, "/") {
				path = path[:len(path)-1]
			}
		}

		req.URL.Path = path

		// This is very important to avoid problems with ssl verification for the HOST header
		if !proxy.PreserveHostHeader {
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

func extractListenPath(listenPath string) string {
	var reg = regexp.MustCompile(`(\/\*.+)`)
	return reg.ReplaceAllString(listenPath, "")
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

// Returns metric name for StatsD in "<request method>.<request path>" format
func getStatsdMetricName(req *http.Request) string {
	return fmt.Sprintf(
		"%s.%s",
		strings.ToLower(req.Method),
		strings.Replace(
			// Double underscores
			strings.Replace(req.URL.Path, "_", "__", -1),
			// and replace dots with single underscore
			".",
			"_",
			-1))
}
