package middleware

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
)

type Logger struct{}

func NewLogger() *Logger {
	return &Logger{}
}

func (m *Logger) Serve(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Debugf("Started %s %s", r.Method, r.URL.Path)
		handler.ServeHTTP(w, r)
	})
}
