package router

import (
	"net/http"

	"github.com/dimfeld/httptreemux"
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
