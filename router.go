package janus

import (
	"net/http"

	"github.com/dimfeld/httptreemux"
)

type HandlerFunc func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)

type Router interface {
	Any(path string, handler http.HandlerFunc)
	Handle(method string, path string, handler http.HandlerFunc)
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
