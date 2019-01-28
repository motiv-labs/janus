package middleware

import (
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/hellofresh/janus/pkg/observability"
)

type reqIDKeyType int

const (
	reqIDKey        reqIDKeyType = iota
	requestIDHeader              = "X-Request-ID"
)

// RequestID middleware
func RequestID(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get(requestIDHeader)
		if requestID == "" {
			requestID = uuid.Must(uuid.NewV4()).String()
		}

		r.Header.Set(requestIDHeader, requestID)
		w.Header().Set(requestIDHeader, requestID)

		handler.ServeHTTP(w, r.WithContext(observability.RequestIDToContext(r.Context(), requestID)))
	})
}
