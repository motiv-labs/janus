package stats

import (
	statsd "gopkg.in/alexcesaro/statsd.v2"
)

type Incrementer struct {
	c *statsd.Client
}

func NewIncrementer(c *statsd.Client) *Incrementer {
	return &Incrementer{c}
}

func (t *Incrementer) Increment(bucket string) {
	t.c.Increment(bucket)
}
