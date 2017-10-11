package proxy

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBalancerFactoryWithSupportedBalance(t *testing.T) {
	balancer, err := NewBalancer("weight")
	assert.NoError(t, err)
	assert.Implements(t, (*Balancer)(nil), balancer)
}

func TestBalancerFactoryWithUnsupportedBalance(t *testing.T) {
	balancer, err := NewBalancer("wrong_alg")
	assert.Error(t, err)
	assert.Nil(t, balancer)
}
