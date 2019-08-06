package middleware

import (
	"net/http"

	"go.opencensus.io/plugin/ochttp/propagation/b3"
	"go.opencensus.io/trace"
	"go.opencensus.io/trace/propagation"
)

const (
	// DebugTraceHeader is the header key used for detecting if
	// trace should be force sampled and returned in the response
	DebugTraceHeader = "X-Debug-Trace"
)

// DebugTrace is a middleware that allows debugging requests by providing the Trace ID
// back to the caller in the same header in the response
func DebugTrace(format propagation.HTTPFormat, key string) func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		if format == nil {
			format = &b3.HTTPFormat{}
		}

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			debugHeader := r.Header.Get(DebugTraceHeader)
			if debugHeader == key {
				ctx, span := trace.StartSpan(r.Context(), DebugTraceHeader, trace.WithSampler(trace.AlwaysSample()))
				r = r.WithContext(ctx)
				format.SpanContextToRequest(span.SpanContext(), r)
				w.Header().Add(DebugTraceHeader, span.SpanContext().TraceID.String())
			}

			handler.ServeHTTP(w, r)
		})
	}
}
