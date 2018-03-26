package middleware

import (
	"context"
	"net/http"

	"github.com/satori/go.uuid"
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
			requestID = uuid.NewV4().String()
		}

		r.Header.Set(requestIDHeader, requestID)
		w.Header().Set(requestIDHeader, requestID)

		ctx := r.Context()
		ctx = context.WithValue(ctx, reqIDKey, requestID)

		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}
