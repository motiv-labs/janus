package stats

import (
	"time"

	log "github.com/Sirupsen/logrus"
	statsd "gopkg.in/alexcesaro/statsd.v2"
)

type MutedTimeTracker struct {
	timer statsd.Timing
	c     *statsd.Client
}

func (t *MutedTimeTracker) Start() {
	t.timer = t.c.NewTiming()
}

func (t *MutedTimeTracker) Finish(bucket string) {
	log.WithFields(log.Fields{
		"bucket":   bucket,
		"elapsed":  int(t.timer.Duration() / time.Millisecond),
		"sampling": "ms",
	}).Debug("Muted stats timer send")
}
