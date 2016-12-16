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

type MiddlewareImplementation interface {
	ProcessRequest(rw http.ResponseWriter, r *http.Request) (error, int)
}

// Generic middleware caller to make extension easier
func CreateMiddleware(mw MiddlewareImplementation) negroni.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		err, errCode := mw.ProcessRequest(rw, r)

		if err != nil {
			c.Abort()
			c.JSON(errCode, err.Error())
		} else {
			next(rw, r)
		}
	}
}
