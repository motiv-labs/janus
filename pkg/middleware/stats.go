package middleware

import (
	"net/http"
	"net/url"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/hellofresh/janus/pkg/response"
	"github.com/hellofresh/stats-go"
)

// Stats represents the stats middleware
type Stats struct {
	statsClient stats.Client
}

// NewStats creates a new instance of Stats
func NewStats(statsClient stats.Client) *Stats {
	return &Stats{statsClient}
}

// Handler is the middleware function
func (m *Stats) Handler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			lock    sync.Mutex
			success bool
		)

		log.WithField("path", r.URL.Path).Debug("Starting Stats middleware")

		hooks := response.Hooks{
			WriteHeader: func(next response.WriteHeaderFunc) response.WriteHeaderFunc {
				return func(code int) {
					next(code)
					lock.Lock()
					defer lock.Unlock()
					if code < http.StatusBadRequest {
						success = true
					}
				}
			},
		}

		timing := m.statsClient.BuildTimer().Start()

		// reverse proxy replaces original request with target request, so keep required fields of the original one
		originalURL := &url.URL{}
		*originalURL = *r.URL
		originalRequest := &http.Request{Method: r.Method, URL: originalURL}

		handler.ServeHTTP(response.Wrap(w, hooks), r)

		log.WithFields(log.Fields{"original_path": originalURL.Path, "request_url": r.URL.Path}).Debug("Track request stats")

		m.statsClient.TrackRequest(originalRequest, timing, success)
	})
}
