package client

import (
	"net/http"
	"sync"
	"time"

	"github.com/hellofresh/stats-go/bucket"
	"github.com/hellofresh/stats-go/incrementer"
	"github.com/hellofresh/stats-go/state"
	"github.com/hellofresh/stats-go/timer"
)

// Metric is a type for storing single duration metric
type Metric struct {
	Bucket  string
	Elapsed time.Duration
}

// Memory is Client implementation for tests
type Memory struct {
	sync.Mutex
	httpMetricCallback bucket.HTTPMetricNameAlterCallback
	httpRequestSection string
	unicode            bool

	TimerMetrics []Metric
	CountMetrics map[string]int
	StateMetrics map[string]int
}

// NewMemory builds and returns new Memory instance
func NewMemory(unicode bool) *Memory {
	client := &Memory{unicode: unicode}
	client.ResetHTTPRequestSection()
	client.resetMetrics()

	return client
}

func (c *Memory) resetMetrics() {
	c.TimerMetrics = []Metric{}
	c.CountMetrics = map[string]int{}
	c.StateMetrics = map[string]int{}
}

// BuildTimer builds timer to track metric timings
func (c *Memory) BuildTimer() timer.Timer {
	return &timer.Memory{}
}

// Close resets all collected stats
func (c *Memory) Close() error {
	c.resetMetrics()
	return nil
}

// TrackRequest tracks HTTP Request stats
func (c *Memory) TrackRequest(r *http.Request, t timer.Timer, success bool) Client {
	b := bucket.NewHTTPRequest(c.httpRequestSection, r, success, c.httpMetricCallback, c.unicode)
	i := incrementer.NewMemory()

	if nil != t {
		c.TimerMetrics = append(c.TimerMetrics, Metric{Bucket: b.Metric(), Elapsed: t.Finish()})
	}

	i.IncrementAll(b)
	for metric, value := range i.Metrics() {
		c.CountMetrics[metric] += value
	}

	return c
}

// TrackOperation tracks custom operation
func (c *Memory) TrackOperation(section string, operation bucket.MetricOperation, t timer.Timer, success bool) Client {
	b := bucket.NewPlain(section, operation, success, true)
	i := incrementer.NewMemory()

	if nil != t {
		c.TimerMetrics = append(c.TimerMetrics, Metric{Bucket: b.MetricWithSuffix(), Elapsed: t.Finish()})
	}

	i.IncrementAll(b)
	for metric, value := range i.Metrics() {
		c.CountMetrics[metric] += value
	}

	return c
}

// TrackOperationN tracks custom operation with n diff
func (c *Memory) TrackOperationN(section string, operation bucket.MetricOperation, t timer.Timer, n int, success bool) Client {
	b := bucket.NewPlain(section, operation, success, true)
	i := incrementer.NewMemory()

	if nil != t {
		c.TimerMetrics = append(c.TimerMetrics, Metric{Bucket: b.MetricWithSuffix(), Elapsed: t.Finish()})
	}

	i.IncrementAllN(b, n)
	for metric, value := range i.Metrics() {
		c.CountMetrics[metric] += value
	}

	return c
}

// TrackMetric tracks custom metric, w/out ok/fail additional sections
func (c *Memory) TrackMetric(section string, operation bucket.MetricOperation) Client {
	b := bucket.NewPlain(section, operation, true, true)
	i := incrementer.NewMemory()

	i.Increment(b.Metric())
	i.Increment(b.MetricTotal())
	for metric, value := range i.Metrics() {
		c.CountMetrics[metric] += value
	}

	return c
}

// TrackMetricN tracks custom metric with n diff, w/out ok/fail additional sections
func (c *Memory) TrackMetricN(section string, operation bucket.MetricOperation, n int) Client {
	b := bucket.NewPlain(section, operation, true, true)
	i := incrementer.NewMemory()

	i.IncrementN(b.Metric(), n)
	i.IncrementN(b.MetricTotal(), n)
	for metric, value := range i.Metrics() {
		c.CountMetrics[metric] += value
	}

	return c
}

// TrackState tracks metric absolute value
func (c *Memory) TrackState(section string, operation bucket.MetricOperation, value int) Client {
	b := bucket.NewPlain(section, operation, true, true)
	s := state.NewMemory()

	s.Set(b.Metric(), value)
	for metric, value := range s.Metrics() {
		c.StateMetrics[metric] = value
	}

	return c
}

// SetHTTPMetricCallback sets callback handler that allows metric operation alteration for HTTP Request
func (c *Memory) SetHTTPMetricCallback(callback bucket.HTTPMetricNameAlterCallback) Client {
	c.Lock()
	defer c.Unlock()

	c.httpMetricCallback = callback
	return c
}

// GetHTTPMetricCallback gets callback handler that allows metric operation alteration for HTTP Request
func (c *Memory) GetHTTPMetricCallback() bucket.HTTPMetricNameAlterCallback {
	c.Lock()
	defer c.Unlock()

	return c.httpMetricCallback
}

// SetHTTPRequestSection sets metric section for HTTP Request metrics
func (c *Memory) SetHTTPRequestSection(section string) Client {
	c.Lock()
	defer c.Unlock()

	c.httpRequestSection = section
	return c
}

// ResetHTTPRequestSection resets metric section for HTTP Request metrics to default value that is "request"
func (c *Memory) ResetHTTPRequestSection() Client {
	return c.SetHTTPRequestSection(bucket.SectionRequest)
}
