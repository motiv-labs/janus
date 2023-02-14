package incrementer

import "github.com/hellofresh/stats-go/bucket"

// Memory struct is Incrementer interface implementation that stores results in memory for further usage
type Memory struct {
	metrics map[string]int
}

// NewMemory builds and returns new Memory instance
func NewMemory() *Memory {
	return &Memory{make(map[string]int)}
}

// Increment increments given metric in memory
func (i *Memory) Increment(metric string) {
	i.metrics[metric]++
}

// IncrementN increments given metric in memory
func (i *Memory) IncrementN(metric string, n int) {
	i.metrics[metric] += n
}

// IncrementAll increments all metrics for given bucket in memory
func (i *Memory) IncrementAll(b bucket.Bucket) {
	incrementAll(i, b)
}

// IncrementAllN increments all metrics for given bucket in memory
func (i *Memory) IncrementAllN(b bucket.Bucket, n int) {
	incrementAllN(i, b, n)
}

// Metrics returns all previously stored metrics
func (i *Memory) Metrics() map[string]int {
	return i.metrics
}
