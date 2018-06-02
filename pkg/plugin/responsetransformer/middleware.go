package responsetransformer

import (
	"net/http"
)

type headerFn func(headerName string, headerValue string)

// Options represents the available options to transform
type Options struct {
	Headers map[string]string `json:"headers"`
}

// Config represent the configuration of the modify headers middleware
type Config struct {
	Add     Options `json:"add"`
	Append  Options `json:"append"`
	Remove  Options `json:"remove"`
	Replace Options `json:"replace"`
}

// NewResponseTransformer creates a new instance of RequestTransformer
func NewResponseTransformer(config Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)

			transform(config.Remove.Headers, removeHeaders(w))
			transform(config.Replace.Headers, replaceHeaders(w))
			transform(config.Add.Headers, addHeaders(w))
			transform(config.Append.Headers, appendHeaders(w))
		})
	}
}

// If and only if the header is not already set, set a new header with the given value. Ignored if the header is already set.
func addHeaders(w http.ResponseWriter) headerFn {
	return func(headerName string, headerValue string) {
		if w.Header().Get(headerName) == "" {
			w.Header().Add(headerName, headerValue)
		}
	}
}

// If the header is not set, set it with the given value. If it is already set, a new header with the same name and the new value will be set.
func appendHeaders(w http.ResponseWriter) headerFn {
	return func(headerName string, headerValue string) {
		w.Header().Add(headerName, headerValue)
	}
}

// Unset the headers with the given name.
func removeHeaders(w http.ResponseWriter) headerFn {
	return func(headerName string, headerValue string) {
		w.Header().Del(headerName)
	}
}

// If and only if the header is already set, replace its old value with the new one. Ignored if the header is not already set.
func replaceHeaders(w http.ResponseWriter) headerFn {
	return func(headerName string, headerValue string) {
		if w.Header().Get(headerName) != "" {
			w.Header().Set(headerName, headerValue)
		}
	}
}

func transform(values map[string]string, fn headerFn) {
	if len(values) <= 0 {
		return
	}

	for name, value := range values {
		fn(name, value)
	}
}
