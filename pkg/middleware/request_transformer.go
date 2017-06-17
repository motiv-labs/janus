package middleware

import (
	"net/http"
	"net/url"
)

type headerFn func(headerName string, headerValue string)

// RequestTransformerOptions represents the available options to transform
type RequestTransformerOptions struct {
	Headers     map[string]string `json:"headers"`
	QueryString map[string]string `json:"querystring"`
}

// RequestTransformerConfig represent the configuration of the modify headers middleware
type RequestTransformerConfig struct {
	Add     RequestTransformerOptions `json:"add"`
	Append  RequestTransformerOptions `json:"append"`
	Remove  RequestTransformerOptions `json:"remove"`
	Replace RequestTransformerOptions `json:"replace"`
}

// RequestTransformer is a middleware that adds or removes headers that go to the upstream
type RequestTransformer struct {
	config RequestTransformerConfig
}

// NewRequestTransformer creates a new instance of RequestTransformer
func NewRequestTransformer(headers RequestTransformerConfig) *RequestTransformer {
	return &RequestTransformer{headers}
}

// Handler is the middleware function
func (h *RequestTransformer) Handler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		h.trasnform(h.config.Remove.Headers, h.removeHeaders(r))
		h.trasnform(h.config.Remove.QueryString, h.removeQueryString(query))

		h.trasnform(h.config.Replace.Headers, h.replaceHeaders(r))
		h.trasnform(h.config.Replace.QueryString, h.replaceQueryString(query))

		h.trasnform(h.config.Add.Headers, h.addHeaders(r))
		h.trasnform(h.config.Add.QueryString, h.addQueryString(query))

		h.trasnform(h.config.Append.Headers, h.appendHeaders(r))
		h.trasnform(h.config.Append.QueryString, h.appendQueryString(query))

		r.URL.RawQuery = query.Encode()

		handler.ServeHTTP(w, r)
	})
}

// If and only if the header is not already set, set a new header with the given value. Ignored if the header is already set.
func (h *RequestTransformer) addHeaders(r *http.Request) headerFn {
	return func(headerName string, headerValue string) {
		if r.Header.Get(headerName) == "" {
			r.Header.Add(headerName, headerValue)
		}
	}
}

// If the header is not set, set it with the given value. If it is already set, a new header with the same name and the new value will be set.
func (h *RequestTransformer) appendHeaders(r *http.Request) headerFn {
	return func(headerName string, headerValue string) {
		r.Header.Add(headerName, headerValue)
	}
}

// Unset the headers with the given name.
func (h *RequestTransformer) removeHeaders(r *http.Request) headerFn {
	return func(headerName string, headerValue string) {
		r.Header.Del(headerName)
	}
}

// If and only if the header is already set, replace its old value with the new one. Ignored if the header is not already set.
func (h *RequestTransformer) replaceHeaders(r *http.Request) headerFn {
	return func(headerName string, headerValue string) {
		if r.Header.Get(headerName) != "" {
			r.Header.Set(headerName, headerValue)
		}
	}
}

func (h *RequestTransformer) trasnform(values map[string]string, fn headerFn) {
	if len(values) <= 0 {
		return
	}

	for name, value := range values {
		fn(name, value)
	}
}

func (h *RequestTransformer) addQueryString(query url.Values) headerFn {
	return func(name string, value string) {
		if query.Get(name) == "" {
			query.Add(name, value)
		}
	}
}

func (h *RequestTransformer) appendQueryString(query url.Values) headerFn {
	return func(name string, value string) {
		query.Add(name, value)
	}
}

func (h *RequestTransformer) removeQueryString(query url.Values) headerFn {
	return func(name string, value string) {
		query.Del(name)
	}
}

func (h *RequestTransformer) replaceQueryString(query url.Values) headerFn {
	return func(name string, value string) {
		if query.Get(name) != "" {
			query.Set(name, value)
		}
	}
}
