package janus

import "net/http"

// Middleware wraps up the APIDefinition object to be included in a
// middleware handler, this can probably be handled better.
type Middleware struct {
	Spec *APISpec
}

// MiddlewareImplementation is an interface that defines how middleware should be implemented.
type MiddlewareImplementation interface {
	ProcessRequest(rw http.ResponseWriter, r *http.Request) (int, error)
}

// CreateMiddleware is a generic middleware caller to make extension easier.
func CreateMiddleware(mw MiddlewareImplementation) HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		statusCode, err := mw.ProcessRequest(rw, r)

		if err != nil {
			// Abort the Request
			rw.WriteHeader(statusCode)
			panic(err)
		} else {
			next(rw, r)
		}
	}
}
