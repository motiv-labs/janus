package middleware

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"

	base "github.com/hellofresh/janus/pkg/opentracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	log "github.com/sirupsen/logrus"
)

var healthChecks = [...]string{
	"/",
	"/status",
	"/metrics",
	"/health",
	"/healthz",
	"/health_check",
	"/healthcheck",
}

type statusCodeTracker struct {
	http.ResponseWriter
	status int
}

func (w *statusCodeTracker) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

// OpenTracing is a middleware that traces the request latency
type OpenTracing struct {
	https bool
}

// NewOpenTracing creates a new instance of OpenTracing
func NewOpenTracing(https bool) *OpenTracing {
	return &OpenTracing{https}
}

// Handler is the middleware function for OpenTracing
func (h *OpenTracing) Handler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isHealthCheck(r.URL.Path) {
			handler.ServeHTTP(w, r)
			return
		}

		var span opentracing.Span
		var err error

		// Attempt to join a trace by getting trace context from the headers.
		wireContext, err := opentracing.GlobalTracer().Extract(
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(r.Header))
		if err != nil {
			// If for whatever reason we can't join, go ahead an start a new root span.
			span = opentracing.StartSpan(r.URL.Path)
		} else {
			span = opentracing.StartSpan(r.URL.Path, opentracing.ChildOf(wireContext))
		}
		defer span.Finish()

		host, err := os.Hostname()
		if host == "" || err != nil {
			log.WithError(err).Debug("Failed to get host name")
			host = "unknown"
		}

		ext.HTTPMethod.Set(span, r.Method)
		ext.HTTPUrl.Set(
			span,
			fmt.Sprintf("%s://%s%s", map[bool]string{true: "https", false: "http"}[h.https], r.Host, r.URL.Path),
		)
		ext.Component.Set(span, "janus")
		ext.SpanKind.Set(span, "server")

		span.SetTag("user.agent", r.UserAgent())
		span.SetTag("peer.address", r.RemoteAddr)
		span.SetTag("host.name", host)
		span.SetTag("referer", r.Referer())
		span.SetTag("request.id", RequestIDFromContext(r.Context()))

		// Add information on the peer service we're about to contact.
		if host, portString, err := net.SplitHostPort(r.URL.Host); err == nil {
			ext.PeerHostname.Set(span, host)
			if port, err := strconv.Atoi(portString); err != nil {
				ext.PeerPort.Set(span, uint16(port))
			}
		} else {
			ext.PeerHostname.Set(span, r.URL.Host)
		}

		err = span.Tracer().Inject(
			span.Context(),
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(r.Header))
		if err != nil {
			log.WithError(err).Error("Could not inject span context into header")
		}

		w = &statusCodeTracker{w, http.StatusOK}
		handler.ServeHTTP(w, base.ToContext(r, span))
		code := uint16(w.(*statusCodeTracker).status)
		if code > http.StatusInternalServerError {
			ext.HTTPStatusCode.Set(span, code)
		}
	})
}

func isHealthCheck(path string) bool {
	for _, hc := range healthChecks {
		if hc == path {
			return true
		}
	}
	return false
}
