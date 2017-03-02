package stats

import (
	statsd "gopkg.in/alexcesaro/statsd.v2"
)

type LiveTimeTracker struct {
	timer statsd.Timing
	c     *statsd.Client
}

func (t *LiveTimeTracker) Start() {
	t.timer = t.c.NewTiming()
}

func (t *LiveTimeTracker) Finish(bucket string) {
	t.timer.Send(bucket)
}
