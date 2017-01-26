package stats

import (
	"fmt"
	"net/http"
	"strings"

	statsd "gopkg.in/alexcesaro/statsd.v2"
)

type StatsClient struct {
	StatsDClient *statsd.Client
}

// NewStatsClient returns initialised stats client instance
func NewStatsClient(statsdClient *statsd.Client) *StatsClient {
	return &StatsClient{statsdClient}
}

// TrackRequest tracks stats for generic request
func (sc *StatsClient) TrackRequest(timing statsd.Timing, req *http.Request) {
	bucket := "request." + sc.getStatsdMetricName(req)

	timing.Send(bucket)
	sc.StatsDClient.Increment(bucket)
}

// TrackRoundTrip tracks stats for round trip request
func (sc *StatsClient) TrackRoundTrip(timing statsd.Timing, req *http.Request, success bool) {
	prefix := fmt.Sprintf("round-%s.", map[bool]string{true: "ok", false: "fail"}[success])
	bucket := prefix + sc.getStatsdMetricName(req)

	timing.Send(bucket)
	sc.StatsDClient.Increment(bucket)
}

// Returns metric name for StatsD in "<request method>.<request path>" format
func (sc *StatsClient) getStatsdMetricName(req *http.Request) string {
	path := strings.Replace(
		// Double underscores
		strings.Replace(req.URL.Path, "_", "__", -1),
		// and replace dots with single underscore
		".",
		"_",
		-1,
	)
	return fmt.Sprintf("%s.%s", strings.ToLower(req.Method), path)
}
