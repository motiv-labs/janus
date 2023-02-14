package bucket

import (
	"net/http"
	"strings"
)

// HTTPMetricNameAlterCallback is a type for HTTP Request metric alter handler
type HTTPMetricNameAlterCallback func(metricParts MetricOperation, r *http.Request) MetricOperation

// HTTPRequest struct in an implementation of Bucket interface that produces metric names for HTTP Request.
// Metrics has the following formats for methods:
//  Metric() -> <section>.<method>.<path-level-0>.<path-level-1>
//  MetricWithSuffix() -> <section>-ok|fail.<method>.<path-level-0>.<path-level-1>
//  TotalRequests() -> total.<section>
//  MetricTotalWithSuffix() -> total-ok|fail.<section>
//
// Normally "<section>" is set to "request", but you can use any string value here.
type HTTPRequest struct {
	*Plain

	r        *http.Request
	callback HTTPMetricNameAlterCallback
}

// NewHTTPRequest builds and returns new HTTPRequest instance
func NewHTTPRequest(section string, r *http.Request, success bool, callback HTTPMetricNameAlterCallback, unicode bool) *HTTPRequest {
	operation := BuildHTTPRequestMetricOperation(r, callback)
	return &HTTPRequest{NewPlain(section, operation, success, unicode), r, callback}
}

// BuildHTTPRequestMetricOperation builds metric operation from HTTP request
func BuildHTTPRequestMetricOperation(r *http.Request, callback HTTPMetricNameAlterCallback) MetricOperation {
	metricParts := MetricOperation{strings.ToLower(r.Method), MetricEmptyPlaceholder, MetricEmptyPlaceholder}
	if r.URL.Path != "/" {
		partsFilled := 1
		for _, fragment := range strings.Split(r.URL.Path, "/") {
			if fragment == "" {
				continue
			}

			metricParts[partsFilled] = fragment
			partsFilled++
			if partsFilled >= len(metricParts) {
				break
			}
		}
	}

	if callback != nil {
		metricParts = callback(metricParts, r)
	}

	return metricParts
}
