package main

import (
	"github.com/hellofresh/api-gateway/storage"
	"github.com/kataras/iris"
	"github.com/valyala/fasthttp"
)

// a silly example
type Database struct {
	*Middleware
	dba *storage.DatabaseAccessor
}

//Important staff, iris middleware must implement the iris.Handler interface which is:
func (m Database) ProcessRequest(req fasthttp.Request, resp fasthttp.Response, c *iris.Context) (error, int) {
	m.Logger.Debug("Starting Database middleware")

	reqSession := m.dba.Clone()
	defer reqSession.Close()
	m.dba.Set(c, reqSession)

	return nil, fasthttp.StatusOK
}
