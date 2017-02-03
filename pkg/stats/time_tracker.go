package stats

import statsd "gopkg.in/alexcesaro/statsd.v2"

type TimeTracker struct {
	timer statsd.Timing
	c     *statsd.Client
}

func NewTimeTracker(c *statsd.Client) *TimeTracker {
	return &TimeTracker{c: c}
}

func (t *TimeTracker) Start() {
	t.timer = t.c.NewTiming()
}

func (t *TimeTracker) Finish(bucket string) {
	t.timer.Send(bucket)
}
