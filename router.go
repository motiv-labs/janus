package janus

import (
	"context"
	"net/http"

	"github.com/dimfeld/httptreemux"
	"github.com/urfave/negroni"
)

// ParamsContextKey is used to retrieve a path's params map from a request's context.
const ParamsContextKey = "params.context.key"

type Params struct {
	params map[string]string
}

// ContextParams returns the params map associated with the given context if one exists. Otherwise, an empty map is returned.
func FromContext(ctx context.Context) *Params {
	if p, ok := ctx.Value(ParamsContextKey).(map[string]string); ok {
		return &Params{p}
	}
	return &Params{}
}

func (p *Params) ByName(name string) string {
	return p.params[name]
}

type HandlerFunc func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)

type Router interface {
	Any(path string, handler http.HandlerFunc)
	Handle(method string, path string, handler http.HandlerFunc)
	GET(path string, handler http.HandlerFunc)
	POST(path string, handler http.HandlerFunc)
	PUT(path string, handler http.HandlerFunc)
	DELETE(path string, handler http.HandlerFunc)
	PATCH(path string, handler http.HandlerFunc)
	HEAD(path string, handler http.HandlerFunc)
	OPTIONS(path string, handler http.HandlerFunc)
	Group(path string) Router
	Use(handler HandlerFunc)
}

type HttpTreeMuxRouter struct {
	innerRouter *httptreemux.ContextGroup
	negroni     *negroni.Negroni
}

func NewHttpTreeMuxRouter() Router {
	router := httptreemux.New()
	return &HttpTreeMuxRouter{
		router.UsingContext(),
		negroni.New(),
	}
}

func (r *HttpTreeMuxRouter) Any(path string, handler http.HandlerFunc) {
	methods := []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodConnect,
		http.MethodHead,
		http.MethodOptions,
		http.MethodTrace,
	}

	for _, method := range methods {
		r.Handle(method, path, handler)
	}
}

func (r *HttpTreeMuxRouter) Handle(method string, path string, handler http.HandlerFunc) {
	r.innerRouter.Handle(method, path, handler)
}

func (r *HttpTreeMuxRouter) GET(path string, handler http.HandlerFunc) {
	r.negroni.UseHandlerFunc(handler)
	r.innerRouter.GET(path, handler)
}

func (r *HttpTreeMuxRouter) POST(path string, handler http.HandlerFunc) {
	r.negroni.UseHandlerFunc(handler)
	r.innerRouter.POST(path, handler)
}

func (r *HttpTreeMuxRouter) PUT(path string, handler http.HandlerFunc) {
	r.negroni.UseHandlerFunc(handler)
	r.innerRouter.PUT(path, handler)
}

func (r *HttpTreeMuxRouter) DELETE(path string, handler http.HandlerFunc) {
	r.negroni.UseHandlerFunc(handler)
	r.innerRouter.DELETE(path, handler)
}

func (r *HttpTreeMuxRouter) PATCH(path string, handler http.HandlerFunc) {
	r.negroni.UseHandlerFunc(handler)
	r.innerRouter.PATCH(path, handler)
}

func (r *HttpTreeMuxRouter) HEAD(path string, handler http.HandlerFunc) {
	r.negroni.UseHandlerFunc(handler)
	r.innerRouter.HEAD(path, handler)
}

func (r *HttpTreeMuxRouter) OPTIONS(path string, handler http.HandlerFunc) {
	r.negroni.UseHandlerFunc(handler)
	r.innerRouter.OPTIONS(path, handler)
}

func (r *HttpTreeMuxRouter) Group(path string) Router {
	return &HttpTreeMuxRouter{
		r.innerRouter.NewContextGroup(path),
		negroni.New(),
	}
}

func (r *HttpTreeMuxRouter) Use(handler HandlerFunc) {
	r.negroni.UseFunc(handler)
}
