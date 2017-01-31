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
		log.Debug("Starting Stats middleware")

		timing := m.statsClient.StatsDClient.NewTiming()

		// reverse proxy replaces original request with target request, so keep required fields of the original one
		originalURL := &url.URL{}
		*originalURL = *r.URL
		originalRequest := &http.Request{Method: r.Method, URL: originalURL}

		handler.ServeHTTP(w, r)

		m.statsClient.TrackRequest(timing, originalRequest)
	})
}
