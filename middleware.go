package janus

import (
	"net/http"

	"github.com/urfave/negroni"
)

// Middleware wraps up the APIDefinition object to be included in a
// middleware handler, this can probably be handled better.
type Middleware struct {
	Spec *APISpec
}

// MiddlewareImplementation is an interface that defines how middleware should be implemented.
type MiddlewareImplementation interface {
	ProcessRequest(r *http.Request, rw http.ResponseWriter) (int, error)
}

// CreateMiddleware is a generic middleware caller to make extension easier.
func CreateMiddleware(mw MiddlewareImplementation) negroni.HandlerFunc {
	return func(r *http.Request, rw http.ResponseWriter, next http.HandlerFunc) {
		statusCode, err := mw.ProcessRequest(r, rw)

		if err != nil {
			// Abort the Request
			rw.WriteHeader(statusCode)
			panic(err)
		} else {
			next(r, rw)
		}
	}
}
