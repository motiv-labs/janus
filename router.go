package janus

import (
	"net/http"

	"github.com/dimfeld/httptreemux"
)

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
}

type HttpTreeMuxRouter struct {
	*httptreemux.ContextGroup
}

func NewHttpTreeMuxRouter() Router {
	router := httptreemux.New()
	return &HttpTreeMuxRouter{
		router.UsingContext(),
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

func (r *HttpTreeMuxRouter) Group(path string) Router {
	return &HttpTreeMuxRouter{r.NewContextGroup(path)}
}
