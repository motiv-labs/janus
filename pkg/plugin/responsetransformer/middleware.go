package responsetransformer

import (
	"net/http"

	"github.com/hellofresh/janus/pkg/proxy"
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
func NewResponseTransformer(config Config) proxy.OutLink {
	return func(req *http.Request, res *http.Response) (*http.Response, error) {
		transform(config.Remove.Headers, removeHeaders(res))
		transform(config.Replace.Headers, replaceHeaders(res))
		transform(config.Add.Headers, addHeaders(res))
		transform(config.Append.Headers, appendHeaders(res))

		return res, nil
	}
}

// If and only if the header is not already set, set a new header with the given value. Ignored if the header is already set.
func addHeaders(res *http.Response) headerFn {
	return func(headerName string, headerValue string) {
		if res.Header.Get(headerName) == "" {
			res.Header.Add(headerName, headerValue)
		}
	}
}

// If the header is not set, set it with the given value. If it is already set, a new header with the same name and the new value will be set.
func appendHeaders(res *http.Response) headerFn {
	return func(headerName string, headerValue string) {
		res.Header.Add(headerName, headerValue)
	}
}

// Unset the headers with the given name.
func removeHeaders(res *http.Response) headerFn {
	return func(headerName string, headerValue string) {
		res.Header.Del(headerName)
	}
}

// If and only if the header is already set, replace its old value with the new one. Ignored if the header is not already set.
func replaceHeaders(res *http.Response) headerFn {
	return func(headerName string, headerValue string) {
		if res.Header.Get(headerName) != "" {
			res.Header.Set(headerName, headerValue)
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
