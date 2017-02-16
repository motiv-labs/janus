package stats

import (
	statsd "gopkg.in/alexcesaro/statsd.v2"
)

type Incrementer interface {
	Increment(bucket string)
}

func NewIncrementer(c *statsd.Client, muted bool) Incrementer {
	if muted {
		return &MutedIncrementer{}
	} else {
		return &LiveIncrementer{c}
	}
}
