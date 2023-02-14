package state

// Memory struct is State interface implementation that stores results in memory for further usage
type Memory struct {
	metrics map[string]int
}

// NewMemory builds and returns new Memory instance
func NewMemory() *Memory {
	return &Memory{make(map[string]int)}
}

// Set sets metric state
func (i *Memory) Set(metric string, n int) {
	i.metrics[metric] = n
}

// Metrics returns all previously stored metrics
func (i *Memory) Metrics() map[string]int {
	return i.metrics
}
