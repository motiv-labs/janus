package proxy

import (
	"errors"
	"reflect"
)

var (
	// ErrUnsupportedAlgorithm is used when an unsupported algorithm is given
	ErrUnsupportedAlgorithm = errors.New("unsupported balancing algorithm")
	typeRegistry            = make(map[string]reflect.Type)
)

func init() {
	typeRegistry["roundrobin"] = reflect.TypeOf(RoundrobinBalancer{})
	typeRegistry["weight"] = reflect.TypeOf(WeightBalancer{})
}

//NewBalancer creates a new Balancer based on balancing strategy
func NewBalancer(balance string) (Balancer, error) {
	alg, ok := typeRegistry[balance]
	if !ok {
		return nil, ErrUnsupportedAlgorithm
	}

	return reflect.New(alg).Elem().Addr().Interface().(Balancer), nil
}
