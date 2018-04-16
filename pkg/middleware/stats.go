package middleware

import (
	"net/http"
	"net/url"

	"github.com/felixge/httpsnoop"
	"github.com/hellofresh/janus/pkg/metrics"
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
func (s *Stats) Handler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithField("path", r.URL.Path).Debug("Starting Stats middleware")
		r = r.WithContext(metrics.NewContext(r.Context(), s.statsClient))

		timing := s.statsClient.BuildTimer().Start()

		// reverse proxy replaces original request with target request, so keep required fields of the original one
		originalURL := &url.URL{}
		*originalURL = *r.URL
		originalRequest := &http.Request{Method: r.Method, URL: originalURL}

		m := httpsnoop.CaptureMetrics(handler, w, r)

		log.WithFields(log.Fields{
			"original_path": originalURL.Path,
			"request_url":   r.URL.Path,
		}).Debug("Track request stats")

		success := m.Code < http.StatusBadRequest
		if m.Code == http.StatusNotFound {
			log.WithField("path", originalURL.Path).Warn("Unknown endpoint requested")
			originalURL.Path = notFoundPath
		}
		s.statsClient.TrackRequest(originalRequest, timing, success)
	})
}
