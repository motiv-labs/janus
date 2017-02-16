package stats

import (
	statsd "gopkg.in/alexcesaro/statsd.v2"
)

type TimeTracker interface {
	Start()
	Finish(bucket string)
}

func NewTimeTracker(c *statsd.Client, muted bool) TimeTracker {
	if muted {
		return &MutedTimeTracker{c: c}
	} else {
		return &LiveTimeTracker{c: c}
	}
}
