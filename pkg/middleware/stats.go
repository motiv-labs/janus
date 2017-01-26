package middleware

import (
	"net/http"

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

		handler.ServeHTTP(w, r)

		m.statsClient.TrackRequest(timing, r)
	})
}
