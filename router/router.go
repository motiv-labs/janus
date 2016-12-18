package router

import (
	"context"
	"net/http"

	"github.com/dimfeld/httptreemux"
	"github.com/justinas/alice"
)

// ParamsContextKey is used to retrieve a path's params map from a request's context.
const ParamsContextKey = "params.context.key"

type MiddlewareImp interface {
	Serve(handler http.Handler) http.Handler
}

type Params struct {
	params map[string]string
}

// FromContext returns the params map associated with the given context if one exists. Otherwise, an empty map is returned.
func FromContext(ctx context.Context) *Params {
	if p, ok := ctx.Value(ParamsContextKey).(map[string]string); ok {
		return &Params{p}
	}
	return &Params{}
}

func (p *Params) ByName(name string) string {
	return p.params[name]
}

type Router interface {
	Handle(method string, path string, handler http.HandlerFunc, middleware ...MiddlewareImp)
	Any(path string, handler http.HandlerFunc, middleware ...MiddlewareImp)
	GET(path string, handler http.HandlerFunc, middleware ...MiddlewareImp)
	POST(path string, handler http.HandlerFunc, middleware ...MiddlewareImp)
	PUT(path string, handler http.HandlerFunc, middleware ...MiddlewareImp)
	DELETE(path string, handler http.HandlerFunc, middleware ...MiddlewareImp)
	PATCH(path string, handler http.HandlerFunc, middleware ...MiddlewareImp)
	HEAD(path string, handler http.HandlerFunc, middleware ...MiddlewareImp)
	OPTIONS(path string, handler http.HandlerFunc, middleware ...MiddlewareImp)
	Group(path string) Router
	Use(handlers ...MiddlewareImp) Router
	Dump() string
}

type HttpTreeMuxRouter struct {
	mux         *httptreemux.TreeMux
	innerRouter *httptreemux.ContextGroup
	chain       alice.Chain
}

func NewHttpTreeMuxRouter() *HttpTreeMuxRouter {
	router := httptreemux.New()
	return &HttpTreeMuxRouter{
		mux:         router,
		innerRouter: router.UsingContext(),
		chain:       alice.New(),
	}
}

func (r *HttpTreeMuxRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

func (r *HttpTreeMuxRouter) Any(path string, handler http.HandlerFunc, middleware ...MiddlewareImp) {
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
		r.Handle(method, path, handler, middleware...)
	}
}

func (r *HttpTreeMuxRouter) Handle(method string, path string, handler http.HandlerFunc, middleware ...MiddlewareImp) {
	var chain alice.Chain
	chain = r.chain

	for _, mw := range middleware {
		chain = chain.Append(mw.Serve)
	}

	r.innerRouter.Handle(method, path, chain.ThenFunc(handler).ServeHTTP)
}

func (r *HttpTreeMuxRouter) GET(path string, handler http.HandlerFunc, middleware ...MiddlewareImp) {
	r.Handle(http.MethodGet, path, handler, middleware...)
}

func (r *HttpTreeMuxRouter) POST(path string, handler http.HandlerFunc, middleware ...MiddlewareImp) {
	r.Handle(http.MethodPost, path, handler, middleware...)
}

func (r *HttpTreeMuxRouter) PUT(path string, handler http.HandlerFunc, middleware ...MiddlewareImp) {
	r.Handle(http.MethodPut, path, handler, middleware...)
}

func (r *HttpTreeMuxRouter) DELETE(path string, handler http.HandlerFunc, middleware ...MiddlewareImp) {
	r.Handle(http.MethodDelete, path, handler, middleware...)
}

func (r *HttpTreeMuxRouter) PATCH(path string, handler http.HandlerFunc, middleware ...MiddlewareImp) {
	r.Handle(http.MethodPatch, path, handler, middleware...)
}

func (r *HttpTreeMuxRouter) HEAD(path string, handler http.HandlerFunc, middleware ...MiddlewareImp) {
	r.Handle(http.MethodHead, path, handler, middleware...)
}

func (r *HttpTreeMuxRouter) OPTIONS(path string, handler http.HandlerFunc, middleware ...MiddlewareImp) {
	r.Handle(http.MethodOptions, path, handler, middleware...)
}

func (r *HttpTreeMuxRouter) Group(path string) Router {
	return &HttpTreeMuxRouter{
		mux:         r.mux,
		innerRouter: r.innerRouter.NewContextGroup(path),
		chain:       r.chain,
	}
}

func (r *HttpTreeMuxRouter) Use(handlers ...MiddlewareImp) Router {
	for _, mw := range handlers {
		r.chain = r.chain.Append(mw.Serve)
	}
	return r
}

func (r *HttpTreeMuxRouter) Dump() string {
	return r.mux.Dump()
}
