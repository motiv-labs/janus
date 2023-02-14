package timer

import "time"

// Memory struct is Timer interface implementation that stores results in memory for further usage
type Memory struct {
	timerStart time.Time
}

// StartAt starts timer at a given time
func (t *Memory) StartAt(s time.Time) Timer {
	t.timerStart = s
	return t
}

// Start starts timer
func (t *Memory) Start() Timer {
	t.timerStart = time.Now()
	return t
}

// Finish returns elapsed duration
func (t *Memory) Finish() time.Duration {
	return time.Now().Sub(t.timerStart)
}
