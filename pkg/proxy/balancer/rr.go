package balancer

import "sync"

type (
	// RoundrobinBalancer balancer
	RoundrobinBalancer struct {
		current int // current backend position
		mu      sync.RWMutex
	}
)

// NewRoundrobinBalancer creates a new instance of Roundrobin
func NewRoundrobinBalancer() *RoundrobinBalancer {
	return &RoundrobinBalancer{}
}

// Elect backend using roundrobin strategy
func (b *RoundrobinBalancer) Elect(hosts []*Target) (*Target, error) {
	if len(hosts) == 0 {
		return nil, ErrEmptyBackendList
	}

	if b.current >= len(hosts) {
		b.current = 0
	}

	host := hosts[b.current]

	b.mu.Lock()
	defer b.mu.Unlock()
	b.current++

	return host, nil
}
