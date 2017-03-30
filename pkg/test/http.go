package test

import (
	"bytes"
	"net/http"
	"net/http/httptest"

	"github.com/hellofresh/janus/pkg/router"
)

// Server represents a testing HTTP Server
type Server struct {
	*httptest.Server
}

// NewServer creates a new instance of Server
func NewServer(r router.Router) *Server {
	return &Server{httptest.NewServer(r)}
}

// Do creates a HTTP request to be tested
func (s *Server) Do(method string, url string, headers map[string]string) (*http.Response, error) {
	var u bytes.Buffer
	u.WriteString(string(s.URL))
	u.WriteString(url)

	req, err := http.NewRequest(method, u.String(), nil)
	if err != nil {
		return nil, err
	}

	for headerName, headerValue := range headers {
		if headerName == "Host" {
			req.Host = headerValue
		} else {
			req.Header.Set(headerName, headerValue)
		}
	}

	return http.DefaultClient.Do(req)
}

// Record creates a ResponseRecorder for testing
func Record(method string, url string, headers map[string]string, handleFunc http.Handler) (*httptest.ResponseRecorder, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	for headerName, headerValue := range headers {
		if headerName == "Host" {
			req.Host = headerValue
		} else {
			req.Header.Set(headerName, headerValue)
		}
	}

	w := httptest.NewRecorder()
	handleFunc.ServeHTTP(w, req)

	return w, nil
}
