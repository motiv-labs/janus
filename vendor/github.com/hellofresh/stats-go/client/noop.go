package client

import (
	"net/http"
	"sync"

	"github.com/hellofresh/stats-go/bucket"
	"github.com/hellofresh/stats-go/timer"
)

// Noop is Client implementation that does literally nothing
type Noop struct {
	sync.Mutex

	httpMetricCallback bucket.HTTPMetricNameAlterCallback
}

// NewNoop builds and returns new Noop instance
func NewNoop() *Noop {
	return &Noop{}
}

// BuildTimer builds timer to track metric timings
func (c *Noop) BuildTimer() timer.Timer {
	return &timer.Memory{}
}

// Close closes underlying client connection if any
func (c *Noop) Close() error {
	return nil
}

// TrackRequest tracks HTTP Request stats
func (c *Noop) TrackRequest(r *http.Request, t timer.Timer, success bool) Client {
	return c
}

// TrackOperation tracks custom operation
func (c *Noop) TrackOperation(section string, operation bucket.MetricOperation, t timer.Timer, success bool) Client {
	return c
}

// TrackOperationN tracks custom operation with n diff
func (c *Noop) TrackOperationN(section string, operation bucket.MetricOperation, t timer.Timer, n int, success bool) Client {
	return c
}

// TrackMetric tracks custom metric, w/out ok/fail additional sections
func (c *Noop) TrackMetric(section string, operation bucket.MetricOperation) Client {
	return c
}

// TrackMetricN tracks custom metric with n diff, w/out ok/fail additional sections
func (c *Noop) TrackMetricN(section string, operation bucket.MetricOperation, n int) Client {
	return c
}

// TrackState tracks metric absolute value
func (c *Noop) TrackState(section string, operation bucket.MetricOperation, value int) Client {
	return c
}

// SetHTTPMetricCallback sets callback handler that allows metric operation alteration for HTTP Request
func (c *Noop) SetHTTPMetricCallback(callback bucket.HTTPMetricNameAlterCallback) Client {
	c.Lock()
	defer c.Unlock()

	c.httpMetricCallback = callback
	return c
}

// GetHTTPMetricCallback gets callback handler that allows metric operation alteration for HTTP Request
func (c *Noop) GetHTTPMetricCallback() bucket.HTTPMetricNameAlterCallback {
	c.Lock()
	defer c.Unlock()

	return c.httpMetricCallback
}

// SetHTTPRequestSection sets metric section for HTTP Request metrics
func (c *Noop) SetHTTPRequestSection(section string) Client {
	return c
}

// ResetHTTPRequestSection resets metric section for HTTP Request metrics to default value that is "request"
func (c *Noop) ResetHTTPRequestSection() Client {
	return c
}
