package client

import (
	"net/http"
	"sync"
	"time"

	"github.com/hellofresh/stats-go/bucket"
	"github.com/hellofresh/stats-go/incrementer"
	"github.com/hellofresh/stats-go/log"
	"github.com/hellofresh/stats-go/state"
	"github.com/hellofresh/stats-go/timer"
	"gopkg.in/alexcesaro/statsd.v2"
)

// StatsD is Client implementation for statsd
type StatsD struct {
	sync.Mutex
	client             *statsd.Client
	httpMetricCallback bucket.HTTPMetricNameAlterCallback
	httpRequestSection string
	unicode            bool
}

// NewStatsD builds and returns new StatsD instance
func NewStatsD(addr string, prefix string, unicode bool) (*StatsD, error) {
	var options []statsd.Option

	if prefix != "" {
		options = append(options, statsd.Prefix(prefix))
	}

	if addr != "" {
		options = append(options, statsd.Address(addr))
	}

	log.Log("Trying to connect to statsd instance", map[string]interface{}{
		"addr":   addr,
		"prefix": prefix,
	}, nil)

	statsdClient, err := statsd.New(options...)
	if err != nil {
		log.Log("An error occurred while connecting to StatsD", map[string]interface{}{
			"addr":   addr,
			"prefix": prefix,
		}, err)
		return nil, err
	}

	client := &StatsD{client: statsdClient, unicode: unicode}
	client.ResetHTTPRequestSection()

	return client, nil
}

// BuildTimer builds timer to track metric timings
func (c *StatsD) BuildTimer() timer.Timer {
	return &timer.Memory{}
}

// Close statsd connection
func (c *StatsD) Close() error {
	c.client.Close()
	return nil
}

// TrackRequest tracks HTTP Request stats
func (c *StatsD) TrackRequest(r *http.Request, t timer.Timer, success bool) Client {
	b := bucket.NewHTTPRequest(c.httpRequestSection, r, success, c.httpMetricCallback, c.unicode)
	i := incrementer.NewStatsD(c.client)

	if nil != t {
		c.client.Timing(b.Metric(), int(t.Finish()/time.Millisecond))
	}
	i.IncrementAll(b)

	return c
}

// TrackOperation tracks custom operation
func (c *StatsD) TrackOperation(section string, operation bucket.MetricOperation, t timer.Timer, success bool) Client {
	b := bucket.NewPlain(section, operation, success, c.unicode)
	i := incrementer.NewStatsD(c.client)

	if nil != t {
		c.client.Timing(b.MetricWithSuffix(), int(t.Finish()/time.Millisecond))
	}
	i.IncrementAll(b)

	return c
}

// TrackOperationN tracks custom operation with n diff
func (c *StatsD) TrackOperationN(section string, operation bucket.MetricOperation, t timer.Timer, n int, success bool) Client {
	b := bucket.NewPlain(section, operation, success, c.unicode)
	i := incrementer.NewStatsD(c.client)

	if nil != t {
		c.client.Timing(b.MetricWithSuffix(), int(t.Finish()/time.Millisecond))
	}
	i.IncrementAllN(b, n)

	return c
}

// TrackMetric tracks custom metric, w/out ok/fail additional sections
func (c *StatsD) TrackMetric(section string, operation bucket.MetricOperation) Client {
	b := bucket.NewPlain(section, operation, true, c.unicode)
	i := incrementer.NewStatsD(c.client)

	i.Increment(b.Metric())
	i.Increment(b.MetricTotal())

	return c
}

// TrackMetricN tracks custom metric with n diff, w/out ok/fail additional sections
func (c *StatsD) TrackMetricN(section string, operation bucket.MetricOperation, n int) Client {
	b := bucket.NewPlain(section, operation, true, c.unicode)
	i := incrementer.NewStatsD(c.client)

	i.IncrementN(b.Metric(), n)
	i.IncrementN(b.MetricTotal(), n)

	return c
}

// TrackState tracks metric absolute value
func (c *StatsD) TrackState(section string, operation bucket.MetricOperation, value int) Client {
	b := bucket.NewPlain(section, operation, true, c.unicode)
	s := state.NewStatsD(c.client)

	s.Set(b.Metric(), value)

	return c
}

// SetHTTPMetricCallback sets callback handler that allows metric operation alteration for HTTP Request
func (c *StatsD) SetHTTPMetricCallback(callback bucket.HTTPMetricNameAlterCallback) Client {
	c.Lock()
	defer c.Unlock()

	c.httpMetricCallback = callback
	return c
}

// GetHTTPMetricCallback gets callback handler that allows metric operation alteration for HTTP Request
func (c *StatsD) GetHTTPMetricCallback() bucket.HTTPMetricNameAlterCallback {
	c.Lock()
	defer c.Unlock()

	return c.httpMetricCallback
}

// SetHTTPRequestSection sets metric section for HTTP Request metrics
func (c *StatsD) SetHTTPRequestSection(section string) Client {
	c.Lock()
	defer c.Unlock()

	c.httpRequestSection = section
	return c
}

// ResetHTTPRequestSection resets metric section for HTTP Request metrics to default value that is "request"
func (c *StatsD) ResetHTTPRequestSection() Client {
	return c.SetHTTPRequestSection(bucket.SectionRequest)
}
