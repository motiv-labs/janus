package middleware

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	base "github.com/hellofresh/janus/pkg/opentracing"
	opentracing "github.com/opentracing/opentracing-go"
)

// OpenTracing is a middleware that traces the request latency
type OpenTracing struct{}

// NewOpenTracing creates a new instance of OpenTracing
func NewOpenTracing() *OpenTracing {
	return &OpenTracing{}
}

// Handler is the middleware function for OpenTracing
func (h *OpenTracing) Handler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var span opentracing.Span
		var err error

		// Attempt to join a trace by getting trace context from the headers.
		wireContext, err := opentracing.GlobalTracer().Extract(
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(r.Header))
		if err != nil {
			// If for whatever reason we can't join, go ahead an start a new root span.
			span = opentracing.StartSpan(r.RequestURI)
		} else {
			span = opentracing.StartSpan(r.RequestURI, opentracing.ChildOf(wireContext))
		}
		defer span.Finish()

		span.SetTag("component", "janus")
		span.SetTag("http.method", r.Method)
		span.SetTag("http.url", r.RequestURI)
		span.SetTag("peer.address", r.RemoteAddr)
		span.SetTag("span.kind", "server")

		err = span.Tracer().Inject(
			span.Context(),
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(r.Header))
		if err != nil {
			log.WithError(err).Error("Could not inject span context into header")
		}

		handler.ServeHTTP(w, base.ToContext(r, span))
	})
}
