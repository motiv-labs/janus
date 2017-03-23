package router

import (
	"context"
	"net/http"

	"github.com/dimfeld/httptreemux"
	"github.com/justinas/alice"
)

// ParamsContextKey is used to retrieve a path's params map from a request's context.
const ParamsContextKey = "params.context.key"

// Constructor for a piece of middleware.
// Some middleware use this constructor out of the box,
// so in most cases you can just pass somepackage.New
type Constructor func(http.Handler) http.Handler

// Params represents the router params
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

// ByName gets a parameter by name
func (p *Params) ByName(name string) string {
	return p.params[name]
}

// Router defines the basic methods for a router
type Router interface {
	Handle(method string, path string, handler http.HandlerFunc, handlers ...Constructor)
	Any(path string, handler http.HandlerFunc, handlers ...Constructor)
	GET(path string, handler http.HandlerFunc, handlers ...Constructor)
	POST(path string, handler http.HandlerFunc, handlers ...Constructor)
	PUT(path string, handler http.HandlerFunc, handlers ...Constructor)
	DELETE(path string, handler http.HandlerFunc, handlers ...Constructor)
	PATCH(path string, handler http.HandlerFunc, handlers ...Constructor)
	HEAD(path string, handler http.HandlerFunc, handlers ...Constructor)
	OPTIONS(path string, handler http.HandlerFunc, handlers ...Constructor)
	Group(path string) Router
	Use(handlers ...Constructor) Router
}

// HTTPTreeMuxRouter is an adapter for httptreemux that implements the Router interface
type HTTPTreeMuxRouter struct {
	mux         *httptreemux.TreeMux
	innerRouter *httptreemux.ContextGroup
	chain       alice.Chain
}

// Options are the HTTPTreeMuxRouter options
type Options struct {
	NotFoundHandler           http.HandlerFunc
	SafeAddRoutesWhileRunning bool
	RedirectMethodBehavior    map[string]httptreemux.RedirectBehavior
}

// DefaultOptions are the default router options
var DefaultOptions = Options{
	NotFoundHandler:           http.NotFound,
	SafeAddRoutesWhileRunning: true,
	RedirectMethodBehavior: map[string]httptreemux.RedirectBehavior{
		// tree mux router uses Redirect301 behavior by default for paths that differs with slash at the end
		// from registered, that causes problems with some services, e.g. api-v2 OPTIONS /menus/ gives 301 and we
		// want by-pass it as is, but only for OPTIONS method, that is processed by CORS
		http.MethodOptions: httptreemux.UseHandler,
	},
}

// NewHTTPTreeMuxWithOptions creates a new instance of HTTPTreeMuxRouter
// with the provided options
func NewHTTPTreeMuxWithOptions(options Options) *HTTPTreeMuxRouter {
	router := httptreemux.New()
	router.NotFoundHandler = options.NotFoundHandler
	router.SafeAddRoutesWhileRunning = options.SafeAddRoutesWhileRunning
	router.RedirectMethodBehavior = options.RedirectMethodBehavior

	return &HTTPTreeMuxRouter{
		mux:         router,
		innerRouter: router.UsingContext(),
		chain:       alice.New(),
	}
}

// NewHTTPTreeMuxRouter creates a new instance of HTTPTreeMuxRouter
func NewHTTPTreeMuxRouter() *HTTPTreeMuxRouter {
	return NewHTTPTreeMuxWithOptions(DefaultOptions)
}

// ServeHTTP server the HTTP requests
func (r *HTTPTreeMuxRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

// Any register a path to all HTTP methods
func (r *HTTPTreeMuxRouter) Any(path string, handler http.HandlerFunc, handlers ...Constructor) {
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
		r.Handle(method, path, handler, handlers...)
	}
}

// Handle registers a path, method and handlers to the router
func (r *HTTPTreeMuxRouter) Handle(method string, path string, handler http.HandlerFunc, handlers ...Constructor) {
	var chain alice.Chain
	chain = r.chain

	for _, h := range handlers {
		chain = chain.Append(alice.Constructor(h))
	}

	r.innerRouter.Handle(method, path, chain.ThenFunc(handler).ServeHTTP)
}

// GET registers a HTTP GET path
func (r *HTTPTreeMuxRouter) GET(path string, handler http.HandlerFunc, handlers ...Constructor) {
	r.Handle(http.MethodGet, path, handler, handlers...)
}

// POST registers a HTTP POST path
func (r *HTTPTreeMuxRouter) POST(path string, handler http.HandlerFunc, handlers ...Constructor) {
	r.Handle(http.MethodPost, path, handler, handlers...)
}

// PUT registers a HTTP PUT path
func (r *HTTPTreeMuxRouter) PUT(path string, handler http.HandlerFunc, handlers ...Constructor) {
	r.Handle(http.MethodPut, path, handler, handlers...)
}

// DELETE registers a HTTP DELETE path
func (r *HTTPTreeMuxRouter) DELETE(path string, handler http.HandlerFunc, handlers ...Constructor) {
	r.Handle(http.MethodDelete, path, handler, handlers...)
}

// PATCH registers a HTTP PATCH path
func (r *HTTPTreeMuxRouter) PATCH(path string, handler http.HandlerFunc, handlers ...Constructor) {
	r.Handle(http.MethodPatch, path, handler, handlers...)
}

// HEAD registers a HTTP HEAD path
func (r *HTTPTreeMuxRouter) HEAD(path string, handler http.HandlerFunc, middleware ...Constructor) {
	r.Handle(http.MethodHead, path, handler, middleware...)
}

// OPTIONS registers a HTTP OPTIONS path
func (r *HTTPTreeMuxRouter) OPTIONS(path string, handler http.HandlerFunc, middleware ...Constructor) {
	r.Handle(http.MethodOptions, path, handler, middleware...)
}

// Group creates a child router for a specific path
func (r *HTTPTreeMuxRouter) Group(path string) Router {
	return &HTTPTreeMuxRouter{
		mux:         r.mux,
		innerRouter: r.innerRouter.NewContextGroup(path),
		chain:       r.chain,
	}
}

// Use attaches a middleware to the router
func (r *HTTPTreeMuxRouter) Use(handlers ...Constructor) Router {
	for _, h := range handlers {
		r.chain = r.chain.Append(alice.Constructor(h))
	}

	return r
}

func (r *HTTPTreeMuxRouter) wrapConstructor(handlers []Constructor) []alice.Constructor {
	var cons = make([]alice.Constructor, len(handlers))
	for _, m := range handlers {
		cons = append(cons, alice.Constructor(m))
	}
	return cons
}
