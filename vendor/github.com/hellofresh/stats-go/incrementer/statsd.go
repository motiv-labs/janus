package incrementer

import (
	"github.com/hellofresh/stats-go/bucket"
	"gopkg.in/alexcesaro/statsd.v2"
)

// StatsD struct is Incrementer interface implementation that writes all metrics to statsd
type StatsD struct {
	c *statsd.Client
}

// NewStatsD creates new statsd incrementer instance
func NewStatsD(c *statsd.Client) *StatsD {
	return &StatsD{c: c}
}

// Increment increments metric in statsd
func (i *StatsD) Increment(metric string) {
	i.c.Increment(metric)
}

// IncrementN increments metric by n in statsd
func (i *StatsD) IncrementN(metric string, n int) {
	i.c.Count(metric, n)
}

// IncrementAll increments all metrics for given bucket in statsd
func (i *StatsD) IncrementAll(b bucket.Bucket) {
	incrementAll(i, b)
}

// IncrementAllN increments all metrics for given bucket in statsd
func (i *StatsD) IncrementAllN(b bucket.Bucket, n int) {
	incrementAllN(i, b, n)
}
