package main

import (
	"github.com/valyala/fasthttp"
	"github.com/kataras/iris"
	"net/http"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

func ToHandler(handler interface{}) iris.Handler {
	h := fasthttpadaptor.NewFastHTTPHandlerFunc(handler.(http.Handler).ServeHTTP)
	return ToHandlerFastHTTP(h)
}

func ToHandlerFastHTTP(h fasthttp.RequestHandler) iris.Handler {
	return iris.HandlerFunc((func(ctx *iris.Context) {
		h(ctx.RequestCtx)
		ctx.Next()
	}))
}
