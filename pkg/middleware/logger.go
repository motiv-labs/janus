package middleware

import (
	"net/http"
	"net/url"
	"time"

	"github.com/felixge/httpsnoop"
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

		// reverse proxy replaces original request with target request, so keep original one
		originalURL := &url.URL{}
		*originalURL = *r.URL

		fields := log.Fields{
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

		if originalURL.String() != r.URL.String() {
			fields["upstream-host"] = r.URL.Host
			fields["upstream-request"] = r.URL.RequestURI()
		}

		log.WithFields(fields).Info("Completed handling request")
	})
}
