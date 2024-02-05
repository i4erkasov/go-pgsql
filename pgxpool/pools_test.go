package pgxpool

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// LoadBalancingSuite defines the structure for the test suite.
type LoadBalancingTestSuite struct {
	suite.Suite
	pools *Pools
}

// SetupTest is executed before each test and initializes necessary resources.
func (t *LoadBalancingTestSuite) SetupTest() {
	t.pools = &Pools{
		pools: make([]*Pool, 3), // Creating a slice of pools with the necessary length.
		count: 0,
	}
}

// TestLoadBalancing checks the uniformity of load balancing.
func (t *LoadBalancingTestSuite) TestLoadBalancing() {
	indexCounts := make(map[int]int)
	numCalls := 1000
	for i := 0; i < numCalls; i++ {
		index := t.pools.slave(len(t.pools.pools))
		if index != 0 { // Ensuring the master pool is not considered
			indexCounts[index]++
		}
	}

	// As the master pool is not considered, we expect an even distribution among the other pools
	expectedCount := numCalls / (len(t.pools.pools) - 1)
	tolerance := expectedCount / 10 // 10% tolerance

	for index, count := range indexCounts {
		if index == 0 { // Skipping the master pool
			continue
		}
		t.True(count >= expectedCount-tolerance && count <= expectedCount+tolerance,
			"Load balancing is not uniform: expected count around %v, got %v for pool %v", expectedCount, count, index)
	}
}

// TestLoadBalancingSuite runs the test suite.
func TestLoadBalancingSuite(t *testing.T) {
	suite.Run(t, new(LoadBalancingTestSuite))
}
