package middleware

import (
	"net/http"

	"github.com/felixge/httpsnoop"
	"github.com/hellofresh/janus/pkg/metrics"
	"github.com/hellofresh/stats-go/client"
	"github.com/hellofresh/stats-go/timer"
	log "github.com/sirupsen/logrus"
)

const (
	notFoundPath          = "/-not-found-"
	statsSectionRoundTrip = "round"
)

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
		log.WithField("path", r.URL.Path).Debug("Starting Stats middleware")
		r = r.WithContext(metrics.NewContext(r.Context(), m.statsClient))

		mt := httpsnoop.CaptureMetrics(handler, w, r)
		t := timer.NewDuration(mt.Duration)

		success := mt.Code < http.StatusBadRequest
		if mt.Code == http.StatusNotFound {
			log.WithField("path", r.URL.Path).Warn("Unknown endpoint requested")
			r.URL.Path = notFoundPath
		}
		m.statsClient.TrackRequest(r, t, success)

		m.statsClient.SetHTTPRequestSection(statsSectionRoundTrip).
			TrackRequest(r, t, mt.Code < http.StatusInternalServerError).
			ResetHTTPRequestSection()
	})
}
