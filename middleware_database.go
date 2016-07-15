package main

import (
	"github.com/hellofresh/api-gateway/storage"
	"github.com/valyala/fasthttp"
	"github.com/kataras/iris"
)

// a silly example
type Database struct {
	*Middleware
	dba *storage.DatabaseAccessor
}

//Important staff, iris middleware must implement the iris.Handler interface which is:
func (d Database) ProcessRequest(req fasthttp.Request, resp fasthttp.Response, c *iris.Context) (error, int) {
	reqSession := d.dba.Clone()
	defer reqSession.Close()
	d.dba.Set(c, reqSession)

	return nil, fasthttp.StatusOK
}
