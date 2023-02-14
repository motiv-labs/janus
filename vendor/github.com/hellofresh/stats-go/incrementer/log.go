package incrementer

import (
	"github.com/hellofresh/stats-go/bucket"
	"github.com/hellofresh/stats-go/log"
)

// Log struct is Incrementer interface implementation that writes all metrics to log
type Log struct{}

// Increment writes given metric to log
func (i *Log) Increment(metric string) {
	log.Log("Stats counter incremented", map[string]interface{}{
		"metric": metric,
	}, nil)
}

// IncrementN writes given metric to log
func (i *Log) IncrementN(metric string, n int) {
	log.Log("Stats counter incremented by n", map[string]interface{}{
		"metric": metric,
		"n":      n,
	}, nil)
}

// IncrementAll writes all metrics for given bucket to log
func (i *Log) IncrementAll(b bucket.Bucket) {
	incrementAll(i, b)
}

// IncrementAllN writes all metrics for given bucket to log
func (i *Log) IncrementAllN(b bucket.Bucket, n int) {
	incrementAllN(i, b, n)
}
