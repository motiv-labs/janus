package middleware

import (
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
)

type Logger struct {
	debug bool
}

func NewLogger(debug bool) *Logger {
	return &Logger{debug}
}

func (m *Logger) Handler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Debugf("started request %s %s", r.Method, r.URL.Path)

		handler.ServeHTTP(w, r)

		latency := time.Since(start)

		fields := log.Fields{
			"method":     r.Method,
			"request":    r.RequestURI,
			"remote":     r.RemoteAddr,
			"duration":   float64(latency.Nanoseconds()) / float64(1000),
			"referer":    r.Referer(),
			"user-agent": r.UserAgent(),
		}

		if m.debug {
			log.WithFields(fields).Info("completed handling request")
		} else {
			log.Info("completed handling request")
		}
	})
}
