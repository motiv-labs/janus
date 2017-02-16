package stats

import (
	"time"

	log "github.com/Sirupsen/logrus"
	statsd "gopkg.in/alexcesaro/statsd.v2"
)

type TimeTracker struct {
	timer statsd.Timing
	c     *statsd.Client
	muted bool
}

func NewTimeTracker(c *statsd.Client, muted bool) *TimeTracker {
	return &TimeTracker{c: c, muted: muted}
}

func (t *TimeTracker) Start() {
	t.timer = t.c.NewTiming()
}

func (t *TimeTracker) Finish(bucket string) {
	if t.muted {
		log.WithFields(log.Fields{
			"bucket":   bucket,
			"elapsed":  int(t.timer.Duration() / time.Millisecond),
			"sampling": "ms",
		}).Debug("Muted stats timer send")
	} else {
		t.timer.Send(bucket)
	}
}
