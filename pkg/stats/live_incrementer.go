package stats

import (
	statsd "gopkg.in/alexcesaro/statsd.v2"
)

type LiveIncrementer struct {
	c *statsd.Client
}

func (t *LiveIncrementer) Increment(bucket string) {
	t.c.Increment(bucket)
}
