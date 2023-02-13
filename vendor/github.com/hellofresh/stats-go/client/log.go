package client

import (
	"net/http"
	"sync"

	"github.com/hellofresh/stats-go/bucket"
	"github.com/hellofresh/stats-go/incrementer"
	"github.com/hellofresh/stats-go/log"
	"github.com/hellofresh/stats-go/state"
	"github.com/hellofresh/stats-go/timer"
)

// Log is Client implementation for debug log
type Log struct {
	sync.Mutex
	httpMetricCallback bucket.HTTPMetricNameAlterCallback
	httpRequestSection string
	unicode            bool
}

// NewLog builds and returns new Log instance
func NewLog(unicode bool) *Log {
	client := &Log{unicode: unicode}
	client.ResetHTTPRequestSection()

	return client
}

// BuildTimer builds timer to track metric timings
func (c *Log) BuildTimer() timer.Timer {
	return &timer.Memory{}
}

// Close statsd connection
func (c *Log) Close() error {
	return nil
}

// TrackRequest tracks HTTP Request stats
func (c *Log) TrackRequest(r *http.Request, t timer.Timer, success bool) Client {
	b := bucket.NewHTTPRequest(c.httpRequestSection, r, success, c.httpMetricCallback, c.unicode)
	i := &incrementer.Log{}

	if nil != t {
		log.Log("Stats timer finished", map[string]interface{}{
			"bucket":  b.Metric(),
			"elapsed": t.Finish().String(),
		}, nil)
	}

	i.IncrementAll(b)

	return c
}

// TrackOperation tracks custom operation
func (c *Log) TrackOperation(section string, operation bucket.MetricOperation, t timer.Timer, success bool) Client {
	b := bucket.NewPlain(section, operation, success, c.unicode)
	i := &incrementer.Log{}

	if nil != t {
		log.Log("Stats timer finished", map[string]interface{}{
			"bucket":  b.MetricWithSuffix(),
			"elapsed": t.Finish().String(),
		}, nil)
	}
	i.IncrementAll(b)

	return c
}

// TrackOperationN tracks custom operation with n diff
func (c *Log) TrackOperationN(section string, operation bucket.MetricOperation, t timer.Timer, n int, success bool) Client {
	b := bucket.NewPlain(section, operation, success, c.unicode)
	i := &incrementer.Log{}

	if nil != t {
		log.Log("Stats timer finished", map[string]interface{}{
			"bucket":  b.MetricWithSuffix(),
			"elapsed": t.Finish().String(),
		}, nil)
	}
	i.IncrementAllN(b, n)

	return c
}

// TrackMetric tracks custom metric, w/out ok/fail additional sections
func (c *Log) TrackMetric(section string, operation bucket.MetricOperation) Client {
	b := bucket.NewPlain(section, operation, true, c.unicode)
	i := &incrementer.Log{}

	i.Increment(b.Metric())
	i.Increment(b.MetricTotal())

	return c
}

// TrackMetricN tracks custom metric with n diff, w/out ok/fail additional sections
func (c *Log) TrackMetricN(section string, operation bucket.MetricOperation, n int) Client {
	b := bucket.NewPlain(section, operation, true, c.unicode)
	i := &incrementer.Log{}

	i.IncrementN(b.Metric(), n)
	i.IncrementN(b.MetricTotal(), n)

	return c
}

// TrackState tracks metric absolute value
func (c *Log) TrackState(section string, operation bucket.MetricOperation, value int) Client {
	b := bucket.NewPlain(section, operation, true, c.unicode)
	s := &state.Log{}

	s.Set(b.Metric(), value)

	return c
}

// SetHTTPMetricCallback sets callback handler that allows metric operation alteration for HTTP Request
func (c *Log) SetHTTPMetricCallback(callback bucket.HTTPMetricNameAlterCallback) Client {
	c.Lock()
	defer c.Unlock()

	c.httpMetricCallback = callback
	return c
}

// GetHTTPMetricCallback gets callback handler that allows metric operation alteration for HTTP Request
func (c *Log) GetHTTPMetricCallback() bucket.HTTPMetricNameAlterCallback {
	c.Lock()
	defer c.Unlock()

	return c.httpMetricCallback
}

// SetHTTPRequestSection sets metric section for HTTP Request metrics
func (c *Log) SetHTTPRequestSection(section string) Client {
	c.Lock()
	defer c.Unlock()

	c.httpRequestSection = section
	return c
}

// ResetHTTPRequestSection resets metric section for HTTP Request metrics to default value that is "request"
func (c *Log) ResetHTTPRequestSection() Client {
	return c.SetHTTPRequestSection(bucket.SectionRequest)
}
