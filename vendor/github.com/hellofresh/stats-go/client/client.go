package client

import (
	"net/http"

	"github.com/hellofresh/stats-go/bucket"
	"github.com/hellofresh/stats-go/timer"
)

// Client is an interface for different methods of gathering stats
type Client interface {
	// BuildTimer builds timer to track metric timings
	BuildTimer() timer.Timer
	// Close closes underlying client connection if any
	Close() error

	// TrackRequest tracks HTTP Request stats
	TrackRequest(r *http.Request, t timer.Timer, success bool) Client

	// TrackOperation tracks custom operation
	TrackOperation(section string, operation bucket.MetricOperation, t timer.Timer, success bool) Client
	// TrackOperationN tracks custom operation with n diff
	TrackOperationN(section string, operation bucket.MetricOperation, t timer.Timer, n int, success bool) Client

	// TrackMetric tracks custom metric, w/out ok/fail additional sections
	TrackMetric(section string, operation bucket.MetricOperation) Client
	// TrackMetricN tracks custom metric with n diff, w/out ok/fail additional sections
	TrackMetricN(section string, operation bucket.MetricOperation, n int) Client

	// TrackState tracks metric absolute value
	TrackState(section string, operation bucket.MetricOperation, value int) Client

	// SetHTTPMetricCallback sets callback handler that allows metric operation alteration for HTTP Request
	SetHTTPMetricCallback(callback bucket.HTTPMetricNameAlterCallback) Client
	// GetHTTPMetricCallback gets callback handler that allows metric operation alteration for HTTP Request
	GetHTTPMetricCallback() bucket.HTTPMetricNameAlterCallback

	// SetHTTPRequestSection sets metric section for HTTP Request metrics
	SetHTTPRequestSection(section string) Client
	// ResetHTTPRequestSection resets metric section for HTTP Request metrics to default value that is "request"
	ResetHTTPRequestSection() Client
}
