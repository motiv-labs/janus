package state

import "gopkg.in/alexcesaro/statsd.v2"

// StatsD struct is State interface implementation that writes all states to statsd gauge
type StatsD struct {
	c *statsd.Client
}

// NewStatsD creates new statsd state instance
func NewStatsD(c *statsd.Client) *StatsD {
	return &StatsD{c: c}
}

// Set sets metric state
func (s *StatsD) Set(metric string, n int) {
	s.c.Gauge(metric, n)
}
