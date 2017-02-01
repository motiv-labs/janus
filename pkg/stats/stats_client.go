package stats

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	statsd "gopkg.in/alexcesaro/statsd.v2"
)

type StatsClient struct {
	StatsDClient *statsd.Client
}

const (
	bucketTotalRequests = "total.requests"
	bucketTotalRound    = "total.round"
)

// NewStatsClient returns initialised stats client instance
func NewStatsClient(statsdClient *statsd.Client) *StatsClient {
	return &StatsClient{statsdClient}
}

// TrackRequest tracks stats for generic request
func (sc *StatsClient) TrackRequest(timing statsd.Timing, req *http.Request) {
	bucket := "request." + sc.getStatsdMetricName(req.Method, req.URL)

	timing.Send(bucket)
	sc.StatsDClient.Increment(bucket)
	sc.StatsDClient.Increment(bucketTotalRequests)
}

// TrackRoundTrip tracks stats for round trip request
func (sc *StatsClient) TrackRoundTrip(timing statsd.Timing, req *http.Request, success bool) {
	okSuffix := map[bool]string{true: "ok", false: "fail"}[success]
	prefix := fmt.Sprintf("round-%s.", okSuffix)
	bucket := prefix + sc.getStatsdMetricName(req.Method, req.URL)

	timing.Send(bucket)
	sc.StatsDClient.Increment(bucket)
	sc.StatsDClient.Increment(bucketTotalRound)
	sc.StatsDClient.Increment(fmt.Sprintf("%s-%s", bucketTotalRound, okSuffix))
}

// Returns metric name for StatsD in "<request method>.<request path>" format
func (sc *StatsClient) getStatsdMetricName(method string, url *url.URL) string {
	path := strings.Replace(
		// Double underscores
		strings.Replace(url.Path, "_", "__", -1),
		// and replace dots with single underscore
		".",
		"_",
		-1,
	)

	metricFragments := []string{"-", "-"}
	if path != "/" {
		fragmentsFilled := 0
		for _, fragment := range strings.Split(path, "/") {
			if fragment == "" {
				continue
			}

			metricFragments[fragmentsFilled] = fragment
			fragmentsFilled++
			if fragmentsFilled >= len(metricFragments) {
				break
			}
		}
	}
	return fmt.Sprintf("%s.%s", strings.ToLower(method), strings.Join(metricFragments, "."))
}
