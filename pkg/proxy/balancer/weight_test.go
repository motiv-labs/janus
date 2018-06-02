package balancer

import (
	"fmt"
	"math"
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

func (suite *WeightBalancerTestSuite) TestWeightBalancerWeight() {
	balancer := NewWeightBalancer()

	totalSteps := 10000
	percentDiffMax := 10

	for _, weights := range []struct {
		weight0 int
		weight1 int
	}{{50, 50}, {80, 20}, {85, 15}, {90, 15}, {20, 80}, {30, 70}, {5, 95}} {
		hosts := []*Target{
			{Target: "127.0.0.1", Weight: weights.weight0},
			{Target: "http://test.com", Weight: weights.weight1},
		}
		shouldElect0 := totalSteps * weights.weight0 / 100
		shouldElect1 := totalSteps * weights.weight1 / 100

		volatility0 := shouldElect0 * percentDiffMax / 100
		volatility1 := shouldElect1 * percentDiffMax / 100

		elected0 := 0
		elected1 := 0
		for i := 0; i < totalSteps; i++ {
			electedHost, err := balancer.Elect(hosts)
			suite.NoError(err)

			if electedHost == hosts[0] {
				elected0++
			} else {
				elected1++
			}
		}

		electedDiff0 := int(math.Abs(float64(elected0 - shouldElect0)))
		suite.True(
			electedDiff0 < volatility0,
			fmt.Sprintf(
				"totalSteps: %d; percentDiffMax: %d; weight0: %d; shouldElect0: %d; elected0: %d; volatility0: %d; electedDiff0: %d",
				totalSteps, percentDiffMax, weights.weight0, shouldElect0, elected0, volatility0, electedDiff0,
			),
		)

		electedDiff1 := int(math.Abs(float64(elected1 - shouldElect1)))
		suite.True(
			electedDiff1 < volatility1,
			fmt.Sprintf(
				"totalSteps: %d; percentDiffMax: %d; weight1: %d; shouldElect1: %d; elected1: %d; volatility1: %d; electedDiff1: %d",
				totalSteps, percentDiffMax, weights.weight1, shouldElect1, elected1, volatility1, electedDiff1,
			),
		)
	}
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestWeightBalancerTestSuiteTestSuite(t *testing.T) {
	suite.Run(t, new(WeightBalancerTestSuite))
}
