package router

import (
	"net/http"

	"github.com/pressly/chi"
)

// ChiRouter is an adapter for chi router that implements the Router interface
type ChiRouter struct {
	mux chi.Router
}

// NewChiRouterWithOptions creates a new instance of ChiRouter
// with the provided options
func NewChiRouterWithOptions(options Options) *ChiRouter {
	router := chi.NewRouter()
	router.NotFound(options.NotFoundHandler)

	return &ChiRouter{
		mux: router,
	}
}

// NewChiRouter creates a new instance of ChiRouter
func NewChiRouter() *ChiRouter {
	return NewChiRouterWithOptions(DefaultOptions)
}

// ServeHTTP server the HTTP requests
func (r *ChiRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

// Any register a path to all HTTP methods
func (r *ChiRouter) Any(path string, handler http.HandlerFunc, handlers ...Constructor) {
	r.with(handlers...).Handle(path, handler)
}

// Handle registers a path, method and handlers to the router
func (r *ChiRouter) Handle(method string, path string, handler http.HandlerFunc, handlers ...Constructor) {
	switch method {
	case http.MethodGet:
		r.GET(path, handler, handlers...)
	case http.MethodPost:
		r.POST(path, handler, handlers...)
	case http.MethodPut:
		r.PUT(path, handler, handlers...)
	case http.MethodPatch:
		r.PATCH(path, handler, handlers...)
	case http.MethodDelete:
		r.DELETE(path, handler, handlers...)
	case http.MethodHead:
		r.HEAD(path, handler, handlers...)
	case http.MethodOptions:
		r.OPTIONS(path, handler, handlers...)
	}
}

// GET registers a HTTP GET path
func (r *ChiRouter) GET(path string, handler http.HandlerFunc, handlers ...Constructor) {
	r.with(handlers...).Get(path, handler)
}

// POST registers a HTTP POST path
func (r *ChiRouter) POST(path string, handler http.HandlerFunc, handlers ...Constructor) {
	r.with(handlers...).Post(path, handler)
}

// PUT registers a HTTP PUT path
func (r *ChiRouter) PUT(path string, handler http.HandlerFunc, handlers ...Constructor) {
	r.with(handlers...).Put(path, handler)
}

// DELETE registers a HTTP DELETE path
func (r *ChiRouter) DELETE(path string, handler http.HandlerFunc, handlers ...Constructor) {
	r.with(handlers...).Delete(path, handler)
}

// PATCH registers a HTTP PATCH path
func (r *ChiRouter) PATCH(path string, handler http.HandlerFunc, handlers ...Constructor) {
	r.with(handlers...).Patch(path, handler)
}

// HEAD registers a HTTP HEAD path
func (r *ChiRouter) HEAD(path string, handler http.HandlerFunc, handlers ...Constructor) {
	r.with(handlers...).Head(path, handler)
}

// OPTIONS registers a HTTP OPTIONS path
func (r *ChiRouter) OPTIONS(path string, handler http.HandlerFunc, handlers ...Constructor) {
	r.with(handlers...).Options(path, handler)
}

// TRACE registers a HTTP TRACE path
func (r *ChiRouter) TRACE(path string, handler http.HandlerFunc, handlers ...Constructor) {
	r.with(handlers...).Trace(path, handler)
}

// CONNECT registers a HTTP CONNECT path
func (r *ChiRouter) CONNECT(path string, handler http.HandlerFunc, handlers ...Constructor) {
	r.with(handlers...).Connect(path, handler)
}

// Group creates a child router for a specific path
func (r *ChiRouter) Group(path string) Router {
	return &ChiRouter{r.mux.Route(path, nil)}
}

// Use attaches a middleware to the router
func (r *ChiRouter) Use(handlers ...Constructor) Router {
	r.mux.Use(r.wrapConstructor(handlers)...)
	return r
}

func (r *ChiRouter) with(handlers ...Constructor) chi.Router {
	return r.mux.With(r.wrapConstructor(handlers)...)
}

func (r *ChiRouter) wrapConstructor(handlers []Constructor) []func(http.Handler) http.Handler {
	var cons = make([]func(http.Handler) http.Handler, 0)
	for _, m := range handlers {
		cons = append(cons, (func(http.Handler) http.Handler)(m))
	}
	return cons
}
