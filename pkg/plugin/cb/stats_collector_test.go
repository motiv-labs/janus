package cb

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/hellofresh/stats-go/client"
)

func TestCollector(t *testing.T) {
	t.Parallel()

	tests := []struct {
		scenario string
		function func(*testing.T)
	}{
		{
			scenario: "when a collector can be created",
			function: testCollectorCreated,
		},
		{
			scenario: "when a collector cannot be created because the metrics client is nil",
			function: testCollectorNotCreated,
		},
		{
			scenario: "when a collector registry is given",
			function: testCollectorRegistry,
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			test.function(t)
		})
	}
}

func testCollectorCreated(t *testing.T) {
	metricsClient := client.NewNoop()
	_, err := NewStatsCollector("test", metricsClient)

	require.NoError(t, err)
}

func testCollectorNotCreated(t *testing.T) {
	_, err := NewStatsCollector("test", nil)

	require.Error(t, err)
}

func testCollectorRegistry(t *testing.T) {
	c := NewCollectorRegistry(client.NewNoop())
	require.NotNil(t, c)

	c = NewCollectorRegistry(nil)
	require.NotNil(t, c)
}
