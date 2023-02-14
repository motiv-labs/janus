package timer

import (
	"time"
)

// Duration struct is Timer interface implementation that writes all timings to statsd
type Duration struct {
	duration time.Duration
}

// NewDuration creates new duration timer instance
func NewDuration(duration time.Duration) *Duration {
	return &Duration{duration: duration}
}

// StartAt starts timer at a given time
func (t *Duration) StartAt(s time.Time) Timer {
	return t
}

// Start starts timer
func (t *Duration) Start() Timer {
	return t
}

// Finish writes elapsed time for metric to statsd
func (t *Duration) Finish() time.Duration {
	return t.duration
}
