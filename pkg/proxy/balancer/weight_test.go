package balancer

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type WeightBalancerTestSuite struct {
	suite.Suite
	hosts []*Target
}

func (suite *WeightBalancerTestSuite) SetupTest() {
	suite.hosts = []*Target{
		{Target: "127.0.0.1", Weight: 5},
		{Target: "http://test.com", Weight: 10},
		{Target: "http://example.com", Weight: 8},
	}
}

func (suite *WeightBalancerTestSuite) TestWeightBalancer() {
	balancer := NewWeightBalancer()

	electedHost, err := balancer.Elect(suite.hosts)
	suite.NoError(err)
	suite.NotNil(electedHost)
}

func (suite *WeightBalancerTestSuite) TestWeightBalancerEmptyList() {
	balancer := NewWeightBalancer()

	_, err := balancer.Elect([]*Target{})
	suite.Error(err)
}

func (suite *WeightBalancerTestSuite) TestWeightBalancerZeroWeight() {
	balancer := NewWeightBalancer()

	_, err := balancer.Elect([]*Target{{Target: "", Weight: 0}})
	suite.Error(err)
}

func (suite *WeightBalancerTestSuite) TestWeightBalancerZeroWeightForOneTarget() {
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
func TestWeightBalancerTestSuiteTestSuite(t *testing.T) {
	suite.Run(t, new(WeightBalancerTestSuite))
}
