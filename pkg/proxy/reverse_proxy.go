package proxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/go-chi/chi"
	"github.com/hellofresh/janus/pkg/observability"
	"github.com/hellofresh/janus/pkg/proxy/balancer"
	"github.com/hellofresh/janus/pkg/router"
	"github.com/hellofresh/stats-go/bucket"
	"github.com/hellofresh/stats-go/client"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"go.opencensus.io/tag"
	"go.opencensus.io/trace"
)

const (
	statsSection = "upstream"
)

// NewBalancedReverseProxy creates a reverse proxy that is load balanced
func NewBalancedReverseProxy(def *Definition, balancer balancer.Balancer, statsClient client.Client) *httputil.ReverseProxy {
	return &httputil.ReverseProxy{
		Director: createDirector(def, balancer, statsClient),
	}
}

func createDirector(proxyDefinition *Definition, balancer balancer.Balancer, statsClient client.Client) func(req *http.Request) {
	paramNameExtractor := router.NewListenPathParamNameExtractor()
	matcher := router.NewListenPathMatcher()

	return func(req *http.Request) {
		upstream, err := balancer.Elect(proxyDefinition.Upstreams.Targets.ToBalancerTargets())
		if err != nil {
			log.WithError(err).Error("Could not elect one upstream")
			return
		}

		targetURL := upstream.Target

		paramNames := paramNameExtractor.Extract(targetURL)
		parametrizedPath, err := applyParameters(req, targetURL, paramNames)
		if err != nil {
			log.WithError(err).Warn("Unable to extract param from request")
		} else {
			targetURL = parametrizedPath
		}

		log.WithField("target", targetURL).Debug("Target upstream elected")

		target, err := url.Parse(targetURL)
		if err != nil {
			log.WithError(err).WithField("upstream_url", targetURL).Error("Could not parse the target URL")
			return
		}

		originalURI := req.RequestURI
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

		log.WithField("path", path).Debug("Upstream Path")
		req.URL.Path = path

		// This is very important to avoid problems with ssl verification for the HOST header
		if proxyDefinition.PreserveHost {
			log.Debug("Preserving the host header")
		} else {
			req.Host = target.Host
		}

		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}

		// Since director modifies cloned request there is no way (or I just did not find one)
		// to get upstream from logger middleware, so we're logging original request and upstream here
		// with the same logging level. Original request is here to match two log messages in case
		// RequestID is not enabled.
		log.WithFields(log.Fields{
			"request":          originalURI,
			"request-id":       observability.RequestIDFromContext(req.Context()),
			"upstream-host":    req.URL.Host,
			"upstream-request": req.URL.RequestURI(),
		}).Info("Proxying request to the following upstream")

		statsClient.TrackMetric(statsSection, bucket.MetricOperation{req.Host})

		// Add additional trace attributes
		addTraceAttributes(req)

		// Insert additional tags
		ctx, _ := tag.New(req.Context(), tag.Insert(observability.KeyUpstreamPath, upstream.Target))
		*req = *req.WithContext(ctx)
	}
}

func addTraceAttributes(req *http.Request) {
	ctx := req.Context()
	span := trace.FromContext(ctx)
	if span == nil {
		return
	}

	host, err := os.Hostname()
	if host == "" || err != nil {
		log.WithError(err).Debug("Failed to get host name")
		host = "unknown"
	}

	span.AddAttributes(
		trace.StringAttribute("http.host", host),
		trace.StringAttribute("http.referrer", req.Referer()),
		trace.StringAttribute("http.remote_address", req.RemoteAddr),
		trace.StringAttribute("request.id", observability.RequestIDFromContext(ctx)),
	)
}

func applyParameters(req *http.Request, path string, paramNames []string) (string, error) {
	for _, paramName := range paramNames {
		paramValue := chi.URLParam(req, paramName)

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
