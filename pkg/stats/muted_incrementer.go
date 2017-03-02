package stats

import (
	log "github.com/Sirupsen/logrus"
)

type MutedIncrementer struct{}

func (t *MutedIncrementer) Increment(bucket string) {
	log.WithField("bucket", bucket).Debug("Muted stats counter increment")
}
