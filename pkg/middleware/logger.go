package middleware

import (
	"net/http"
	"time"

	"github.com/felixge/httpsnoop"
	"github.com/hellofresh/janus/pkg/observability"
	log "github.com/sirupsen/logrus"
)

// Logger struct contains data and logic required for middleware functionality
type Logger struct{}

// NewLogger builds and returns new Logger middleware instance
func NewLogger() *Logger {
	return &Logger{}
}

// Handler implementation
func (m *Logger) Handler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{"method": r.Method, "path": r.URL.Path}).Debug("Started request")

		fields := log.Fields{
			"request-id":  observability.RequestIDFromContext(r.Context()),
			"method":      r.Method,
			"host":        r.Host,
			"request":     r.RequestURI,
			"remote-addr": r.RemoteAddr,
			"referer":     r.Referer(),
			"user-agent":  r.UserAgent(),
		}

		m := httpsnoop.CaptureMetrics(handler, w, r)

		fields["code"] = m.Code
		fields["duration"] = int(m.Duration / time.Millisecond)
		fields["duration-fmt"] = m.Duration.String()

		log.WithFields(fields).Info("Completed handling request")
	})
}
