package stats

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	totalRequestBucket   = "total.requests"
	totalRoundTripBucket = "total.round"
	pathIDPlaceholder    = "-id-"
)

func testAlwaysTrue(string) bool {
	return true
}

func testIsNumeric(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

// key - path first level
// value - function to test if the second level is ID
var hasIDAtSecondLevel = map[string]func(string) bool{
	"user":              testAlwaysTrue,
	"users":             testAlwaysTrue,
	"allergens":         testAlwaysTrue,
	"cuisines":          testAlwaysTrue,
	"favorites":         testAlwaysTrue,
	"ingredients":       testAlwaysTrue,
	"menus":             testAlwaysTrue,
	"ratings":           testAlwaysTrue,
	"recipes":           testAlwaysTrue,
	"addresses":         testAlwaysTrue,
	"boxes":             testAlwaysTrue,
	"coupons":           testAlwaysTrue,
	"customers":         testAlwaysTrue,
	"delivery__options": testAlwaysTrue,
	"product__families": testAlwaysTrue,
	"products":          testAlwaysTrue,
	"subscriptions":     testIsNumeric,
}

type Bucket interface {
	Name() string
}

func RequestsWithSuffixBucket(r *http.Request, success bool) string {
	return fmt.Sprintf("request-%s.%s", getRequestStatus(success), getMetricName(r.Method, r.URL))
}

func TotalRequestsWithSuffixBucket(success bool) string {
	return fmt.Sprintf("%s-%s", totalRequestBucket, getRequestStatus(success))
}

func RequestBucket(r *http.Request) string {
	return fmt.Sprintf("request.%s", getMetricName(r.Method, r.URL))
}

func RoundTripBucket(r *http.Request, success bool) string {
	return fmt.Sprintf("round-%s.%s", getRequestStatus(success), getMetricName(r.Method, r.URL))
}

func RoundTripSuffixBucket(success bool) string {
	return fmt.Sprintf("%s-%s", totalRoundTripBucket, getRequestStatus(success))
}

func getRequestStatus(success bool) string {
	return map[bool]string{true: "ok", false: "fail"}[success]
}

// Returns metric name for StatsD in "<request method>.<request path>" format
func getMetricName(method string, url *url.URL) string {
	path := sanitizePathForMetric(url.Path)

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

	if testFunction, ok := hasIDAtSecondLevel[metricFragments[0]]; ok {
		if testFunction(metricFragments[1]) {
			metricFragments[1] = pathIDPlaceholder
		}
	}

	return fmt.Sprintf("%s.%s", strings.ToLower(method), strings.Join(metricFragments, "."))
}

func sanitizePathForMetric(path string) string {
	return strings.Replace(
		// Double underscores
		strings.Replace(path, "_", "__", -1),
		// and replace dots with single underscore
		".",
		"_",
		-1,
	)
}
