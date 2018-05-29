// Package balancer provides a simple interface to create concrete balancer algorightms that can be used to choose
// an upstream
package balancer

import (
	"errors"
	"math/rand"
	"reflect"
	"time"
)

var (
	// ErrEmptyBackendList is used when the list of beckends is empty
	ErrEmptyBackendList = errors.New("can not elect backend, Backends empty")
	// ErrCannotElectBackend is used a backend cannot be elected
	ErrCannotElectBackend = errors.New("cant elect backend")
	// ErrUnsupportedAlgorithm is used when an unsupported algorithm is given
	ErrUnsupportedAlgorithm = errors.New("unsupported balancing algorithm")
	typeRegistry            = make(map[string]reflect.Type)
)

type (
	// Balancer holds the load balancer methods for many different algorithms
	Balancer interface {
		Elect(hosts []*Target) (*Target, error)
	}

	// Target is an ip address/hostname with a port that identifies an instance of a backend service
	Target struct {
		Target string
		Weight int
	}
)

func init() {
	rand.Seed(time.Now().UnixNano())
	typeRegistry["roundrobin"] = reflect.TypeOf(RoundrobinBalancer{})
	typeRegistry["rr"] = reflect.TypeOf(RoundrobinBalancer{})
	typeRegistry["weight"] = reflect.TypeOf(WeightBalancer{})
}

// New creates a new Balancer based on balancing strategy
func New(balance string) (Balancer, error) {
	alg, ok := typeRegistry[balance]
	if !ok {
		return nil, ErrUnsupportedAlgorithm
	}

	return reflect.New(alg).Elem().Addr().Interface().(Balancer), nil
}
