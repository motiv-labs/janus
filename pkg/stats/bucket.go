package stats

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const (
	TotalRequestBucket   = "total.requests"
	TotalRoundTripBucket = "total.round"
)

type Bucket interface {
	Name() string
}

func RequestsWithSuffixBucket(r *http.Request, success bool) string {
	return fmt.Sprintf("request-%s.%s", getRequestStatus(success), getMetricName(r.Method, r.URL))
}

func TotalRequestsWithSuffixBucket(success bool) string {
	return fmt.Sprintf("%s-%s", TotalRequestBucket, getRequestStatus(success))
}

func RequestBucket(r *http.Request) string {
	return fmt.Sprintf("request.%s", getMetricName(r.Method, r.URL))
}

func RoundTripBucket(r *http.Request, success bool) string {
	return fmt.Sprintf("round-%s.%s", getRequestStatus(success), getMetricName(r.Method, r.URL))
}

func RoundTripSuffixBucket(success bool) string {
	return fmt.Sprintf("%s-%s", TotalRoundTripBucket, getRequestStatus(success))
}

func getRequestStatus(success bool) string {
	return map[bool]string{true: "ok", false: "fail"}[success]
}

// Returns metric name for StatsD in "<request method>.<request path>" format
func getMetricName(method string, url *url.URL) string {
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
