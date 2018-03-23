package proxy

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type BalancerTestSuite struct {
	suite.Suite
	hosts []*Target
}

func (suite *BalancerTestSuite) SetupTest() {
	suite.hosts = []*Target{
		{Target: "127.0.0.1", Weight: 5},
		{Target: "http://test.com", Weight: 10},
		{Target: "http://example.com", Weight: 8},
	}
}

func (suite *BalancerTestSuite) TestRoundRobinBalancerSuccessfulBalance() {
	balancer := NewRoundrobinBalancer()

	electedHost, err := balancer.Elect(suite.hosts)
	suite.NoError(err)
	suite.Equal(suite.hosts[0], electedHost)

	electedHost, err = balancer.Elect(suite.hosts)
	suite.NoError(err)
	suite.Equal(suite.hosts[1], electedHost)

	electedHost, err = balancer.Elect(suite.hosts)
	suite.NoError(err)
	suite.Equal(suite.hosts[2], electedHost)

	electedHost, err = balancer.Elect(suite.hosts)
	suite.NoError(err)
	suite.Equal(suite.hosts[0], electedHost)
}

func (suite *BalancerTestSuite) TestRoundRobinBalancerEmptyList() {
	balancer := NewRoundrobinBalancer()

	_, err := balancer.Elect([]*Target{})
	suite.Error(err)
}

func (suite *BalancerTestSuite) TestWeightBalancer() {
	balancer := NewWeightBalancer()

	electedHost, err := balancer.Elect(suite.hosts)
	suite.NoError(err)
	suite.NotNil(electedHost)
}

func (suite *BalancerTestSuite) TestWeightBalancerEmptyList() {
	balancer := NewWeightBalancer()

	_, err := balancer.Elect([]*Target{})
	suite.Error(err)
}

func (suite *BalancerTestSuite) TestWeightBalancerZeroWeight() {
	balancer := NewWeightBalancer()

	_, err := balancer.Elect([]*Target{{Target: "", Weight: 0}})
	suite.Error(err)
}

func (suite *BalancerTestSuite) TestWeightBalancerZeroWeightForOneTarget() {
	balancer := NewWeightBalancer()

	hosts := []*Target{
		{Target: "127.0.0.1", Weight: 0},
		{Target: "http://test.com", Weight: 100},
	}

	electedHost, err := balancer.Elect(hosts)
	suite.NoError(err)
	suite.Equal(hosts[1], electedHost)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestBalancerTestSuite(t *testing.T) {
	suite.Run(t, new(BalancerTestSuite))
}
