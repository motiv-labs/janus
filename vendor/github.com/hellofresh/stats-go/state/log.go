package state

import (
	"github.com/hellofresh/stats-go/log"
)

// Log struct is State interface implementation that writes all states to log
type Log struct{}

// Set sets metric state
func (s *Log) Set(metric string, n int) {
	log.Log("Stats state set", map[string]interface{}{
		"bucket": metric,
		"state":  n,
	}, nil)
}
