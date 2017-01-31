package middleware

import (
	"net/http"
	"net/url"

	log "github.com/Sirupsen/logrus"
	"github.com/hellofresh/janus/pkg/stats"
)

type Stats struct {
	statsClient *stats.StatsClient
}

func NewStats(statsClient *stats.StatsClient) *Stats {
	return &Stats{statsClient}
}

func (m *Stats) Handler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithField("path", r.URL.Path).Debug("Starting Stats middleware")

		timing := m.statsClient.StatsDClient.NewTiming()

		// reverse proxy replaces original request with target request, so keep required fields of the original one
		originalURL := &url.URL{}
		*originalURL = *r.URL
		originalRequest := &http.Request{Method: r.Method, URL: originalURL}

		handler.ServeHTTP(w, r)

		log.WithFields(log.Fields{"original_path": originalURL.Path, "request_url": r.URL.Path}).Debug("Track request stats")
		m.statsClient.TrackRequest(timing, originalRequest)
	})
}
