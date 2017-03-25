package test

import (
	"bytes"
	"net/http"
	"net/http/httptest"

	"github.com/hellofresh/janus/pkg/router"
)

type Server struct {
	*httptest.Server
}

func CreateServer(r router.Router) *Server {
	return &Server{httptest.NewServer(r)}
}

func (s *Server) Do(method string, url string) (*http.Response, error) {
	var u bytes.Buffer
	u.WriteString(string(s.URL))
	u.WriteString(url)

	req, err := http.NewRequest(method, u.String(), nil)
	if err != nil {
		return nil, err
	}

	return http.DefaultClient.Do(req)
}

func Record(method string, url string, handleFunc http.Handler) (*httptest.ResponseRecorder, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handleFunc.ServeHTTP(w, req)

	return w, nil
}
