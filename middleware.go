package main

import (
	"github.com/kataras/iris"
	"github.com/valyala/fasthttp"
)

// Middleware wraps up the APIDefinition object to be included in a
// middleware handler, this can probably be handled better.
type Middleware struct {
	Spec *APISpec
}

type MiddlewareImplementation interface {
	ProcessRequest(req fasthttp.Request, resp fasthttp.Response, c *iris.Context) (error, int)
}

// Generic middleware caller to make extension easier
func CreateMiddleware(mw MiddlewareImplementation, tykMwSuper *Middleware) {
	irisHandler := func(c *iris.Context) {
		req := c.Request
		res := c.Response

		reqErr, errCode := mw.ProcessRequest(req, res, c)

		if reqErr != nil {
			c.Error(reqErr.Error(), errCode)
			return
		}

		c.SetStatusCode(errCode)
		c.Next()
	}

	iris.UseFunc(irisHandler)
}
