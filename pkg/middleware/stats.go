package middleware

import (
	"net/http"
	"net/url"
	"sync"

	"github.com/hellofresh/janus/pkg/metrics"
	"github.com/hellofresh/janus/pkg/response"
	"github.com/hellofresh/stats-go/client"
	log "github.com/sirupsen/logrus"
)

const notFoundPath = "/-not-found-"

// Stats represents the stats middleware
type Stats struct {
	statsClient client.Client
}

// NewStats creates a new instance of Stats
func NewStats(statsClient client.Client) *Stats {
	return &Stats{statsClient}
}

// Handler is the middleware function
func (m *Stats) Handler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			lock         sync.Mutex
			responseCode int
		)

		log.WithField("path", r.URL.Path).Debug("Starting Stats middleware")
		r.WithContext(metrics.NewContext(r.Context(), m.statsClient))

		hooks := response.Hooks{
			WriteHeader: func(next response.WriteHeaderFunc) response.WriteHeaderFunc {
				return func(code int) {
					next(code)
					lock.Lock()
					defer lock.Unlock()
					responseCode = code
				}
			},
		}

		timing := m.statsClient.BuildTimer().Start()

		// reverse proxy replaces original request with target request, so keep required fields of the original one
		originalURL := &url.URL{}
		*originalURL = *r.URL
		originalRequest := &http.Request{Method: r.Method, URL: originalURL}

		handler.ServeHTTP(response.Wrap(w, hooks), r)

		log.WithFields(log.Fields{
			"original_path": originalURL.Path,
			"request_url":   r.URL.Path,
		}).Debug("Track request stats")

		success := responseCode < http.StatusBadRequest
		if responseCode == http.StatusNotFound {
			log.WithField("path", originalURL.Path).Warn("Unknown endpoint requested")
			originalURL.Path = notFoundPath
		}
		m.statsClient.TrackRequest(originalRequest, timing, success)
	})
}
