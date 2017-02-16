package stats

import (
	log "github.com/Sirupsen/logrus"
	statsd "gopkg.in/alexcesaro/statsd.v2"
)

type Incrementer struct {
	c     *statsd.Client
	muted bool
}

func NewIncrementer(c *statsd.Client, muted bool) *Incrementer {
	return &Incrementer{c, muted}
}

func (t *Incrementer) Increment(bucket string) {
	if t.muted {
		log.WithField("bucket", bucket).Debug("Muted stats counter increment")
	} else {
		t.c.Increment(bucket)
	}
}
