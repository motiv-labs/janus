package requesttransformer

import (
	"net/http"
	"net/url"
)

type headerFn func(headerName string, headerValue string)

// Options represents the available options to transform
type Options struct {
	Headers     map[string]string `json:"headers"`
	QueryString map[string]string `json:"querystring"`
}

// Config represent the configuration of the modify headers middleware
type Config struct {
	Add     Options `json:"add"`
	Append  Options `json:"append"`
	Remove  Options `json:"remove"`
	Replace Options `json:"replace"`
}

// NewRequestTransformer creates a new instance of RequestTransformer
func NewRequestTransformer(config Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			query := r.URL.Query()

			transform(config.Remove.Headers, removeHeaders(r))
			transform(config.Remove.QueryString, removeQueryString(query))

			transform(config.Replace.Headers, replaceHeaders(r))
			transform(config.Replace.QueryString, replaceQueryString(query))

			transform(config.Add.Headers, addHeaders(r))
			transform(config.Add.QueryString, addQueryString(query))

			transform(config.Append.Headers, appendHeaders(r))
			transform(config.Append.QueryString, appendQueryString(query))

			r.URL.RawQuery = query.Encode()

			next.ServeHTTP(w, r)
		})
	}
}

// If and only if the header is not already set, set a new header with the given value. Ignored if the header is already set.
func addHeaders(r *http.Request) headerFn {
	return func(headerName string, headerValue string) {
		if r.Header.Get(headerName) == "" {
			r.Header.Add(headerName, headerValue)
		}
	}
}

// If the header is not set, set it with the given value. If it is already set, a new header with the same name and the new value will be set.
func appendHeaders(r *http.Request) headerFn {
	return func(headerName string, headerValue string) {
		r.Header.Add(headerName, headerValue)
	}
}

// Unset the headers with the given name.
func removeHeaders(r *http.Request) headerFn {
	return func(headerName string, headerValue string) {
		r.Header.Del(headerName)
	}
}

// If and only if the header is already set, replace its old value with the new one. Ignored if the header is not already set.
func replaceHeaders(r *http.Request) headerFn {
	return func(headerName string, headerValue string) {
		if r.Header.Get(headerName) != "" {
			r.Header.Set(headerName, headerValue)
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

func addQueryString(query url.Values) headerFn {
	return func(name string, value string) {
		if query.Get(name) == "" {
			query.Add(name, value)
		}
	}
}

func appendQueryString(query url.Values) headerFn {
	return func(name string, value string) {
		query.Add(name, value)
	}
}

func removeQueryString(query url.Values) headerFn {
	return func(name string, value string) {
		query.Del(name)
	}
}

func replaceQueryString(query url.Values) headerFn {
	return func(name string, value string) {
		if query.Get(name) != "" {
			query.Set(name, value)
		}
	}
}
