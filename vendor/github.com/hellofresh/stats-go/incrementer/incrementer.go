package incrementer

import (
	"github.com/hellofresh/stats-go/bucket"
)

// Incrementer is a metric incrementer interface
type Incrementer interface {
	// Increment increments metric
	Increment(metric string)

	// IncrementN increments metric by n
	IncrementN(metric string, n int)

	// Increment increments all metrics for given bucket
	IncrementAll(b bucket.Bucket)

	// Increment increments all metrics for given bucket by n
	IncrementAllN(b bucket.Bucket, n int)
}

func incrementAll(i Incrementer, b bucket.Bucket) {
	i.Increment(b.Metric())
	i.Increment(b.MetricWithSuffix())
	i.Increment(b.MetricTotal())
	i.Increment(b.MetricTotalWithSuffix())
}

func incrementAllN(i Incrementer, b bucket.Bucket, n int) {
	i.IncrementN(b.Metric(), n)
	i.IncrementN(b.MetricWithSuffix(), n)
	i.IncrementN(b.MetricTotal(), n)
	i.IncrementN(b.MetricTotalWithSuffix(), n)
}
