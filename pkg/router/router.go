package router

import (
	"net/http"

	"github.com/pressly/chi"
)

// Constructor for a piece of middleware.
// Some middleware use this constructor out of the box,
// so in most cases you can just pass somepackage.New
type Constructor func(http.Handler) http.Handler

// URLParam returns the url parameter from a http.Request object.
func URLParam(r *http.Request, key string) string {
	return chi.URLParam(r, key)
}

// Router defines the basic methods for a router
type Router interface {
	ServeHTTP(w http.ResponseWriter, req *http.Request)
	Handle(method string, path string, handler http.HandlerFunc, handlers ...Constructor)
	Any(path string, handler http.HandlerFunc, handlers ...Constructor)
	GET(path string, handler http.HandlerFunc, handlers ...Constructor)
	POST(path string, handler http.HandlerFunc, handlers ...Constructor)
	PUT(path string, handler http.HandlerFunc, handlers ...Constructor)
	DELETE(path string, handler http.HandlerFunc, handlers ...Constructor)
	PATCH(path string, handler http.HandlerFunc, handlers ...Constructor)
	HEAD(path string, handler http.HandlerFunc, handlers ...Constructor)
	OPTIONS(path string, handler http.HandlerFunc, handlers ...Constructor)
	TRACE(path string, handler http.HandlerFunc, handlers ...Constructor)
	CONNECT(path string, handler http.HandlerFunc, handlers ...Constructor)
	Group(path string) Router
	Use(handlers ...Constructor) Router

	RoutesCount() int
}

// Options are the HTTPTreeMuxRouter options
type Options struct {
	NotFoundHandler           http.HandlerFunc
	SafeAddRoutesWhileRunning bool
}

// DefaultOptions are the default router options
var DefaultOptions = Options{
	NotFoundHandler:           http.NotFound,
	SafeAddRoutesWhileRunning: true,
}
