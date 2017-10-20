package proxy

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBalancerFactoryWithSupportedBalance(t *testing.T) {
	balancer, err := NewBalancer("weight")
	require.NoError(t, err)
	assert.Implements(t, (*Balancer)(nil), balancer)
}

func TestBalancerFactoryWithUnsupportedBalance(t *testing.T) {
	balancer, err := NewBalancer("wrong_alg")
	require.Error(t, err)
	assert.Nil(t, balancer)
}
