package proxy

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	hosts = []*Target{
		&Target{Target: "127.0.0.1", Weight: 5},
		&Target{Target: "http://test.com", Weight: 10},
		&Target{Target: "http://example.com", Weight: 8},
	}
)

func TestRoundRobinBalancer(t *testing.T) {
	balancer := NewRoundrobinBalancer()

	electedHost, err := balancer.Elect(hosts)
	assert.NoError(t, err)
	assert.Equal(t, hosts[0], electedHost)

	electedHost, err = balancer.Elect(hosts)
	assert.NoError(t, err)
	assert.Equal(t, hosts[1], electedHost)

	electedHost, err = balancer.Elect(hosts)
	assert.NoError(t, err)
	assert.Equal(t, hosts[2], electedHost)

	electedHost, err = balancer.Elect(hosts)
	assert.NoError(t, err)
	assert.Equal(t, hosts[0], electedHost)
}

func TestRoundRobinBalancerEmptyList(t *testing.T) {
	balancer := NewRoundrobinBalancer()

	_, err := balancer.Elect([]*Target{})
	assert.Error(t, err)
}

func TestWeightBalancer(t *testing.T) {
	balancer := NewWeightBalancer()

	electedHost, err := balancer.Elect(hosts)
	assert.NoError(t, err)
	assert.Equal(t, hosts[1], electedHost)

	electedHost, err = balancer.Elect(hosts)
	assert.NoError(t, err)
	assert.Equal(t, hosts[2], electedHost)

	electedHost, err = balancer.Elect(hosts)
	assert.NoError(t, err)
	assert.Equal(t, hosts[2], electedHost)

	electedHost, err = balancer.Elect(hosts)
	assert.NoError(t, err)
	assert.Equal(t, hosts[0], electedHost)

	electedHost, err = balancer.Elect(hosts)
	assert.NoError(t, err)
	assert.Equal(t, hosts[1], electedHost)
}

func TestWeightBalancerEmptyList(t *testing.T) {
	balancer := NewWeightBalancer()

	_, err := balancer.Elect([]*Target{})
	assert.Error(t, err)
}

func TestWeightBalancerZeroWeight(t *testing.T) {
	balancer := NewWeightBalancer()

	_, err := balancer.Elect([]*Target{&Target{Target: "", Weight: 0}})
	assert.Error(t, err)
}
