package proxy

import (
	"errors"
	"math/rand"
	"sync"
	"time"
)

var (
	// ErrEmptyBackendList is used when the list of beckends is empty
	ErrEmptyBackendList = errors.New("can not elect backend, Backends empty")
	// ErrZeroWeight is used when there a zero value weight was given
	ErrZeroWeight = errors.New("invalid backend, weight 0 given")
	// ErrCannotElectBackend is used a backend cannot be elected
	ErrCannotElectBackend = errors.New("cant elect backend")
)

type (
	// Balancer holds the load balancer methods for many different algorithms
	Balancer interface {
		Elect(hosts []*Target) (*Target, error)
	}

	// RoundrobinBalancer balancer
	RoundrobinBalancer struct {
		current int // current backend position
		mu      sync.RWMutex
	}

	// WeightBalancer balancer
	WeightBalancer struct{}
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

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

// NewWeightBalancer creates a new instance of Roundrobin
func NewWeightBalancer() *WeightBalancer {
	return &WeightBalancer{}
}

// Elect backend using Weight strategy
func (b *WeightBalancer) Elect(hosts []*Target) (*Target, error) {
	if len(hosts) == 0 {
		return nil, ErrEmptyBackendList
	}

	totalWeight := 0
	for _, host := range hosts {
		totalWeight += host.Weight
	}

	if totalWeight <= 0 {
		return nil, ErrZeroWeight
	}

	r := rand.Intn(totalWeight)
	pos := 0

	for _, host := range hosts {
		pos += host.Weight
		if r >= pos {
			continue
		}
		return host, nil
	}

	return nil, ErrCannotElectBackend
}
