package proxy

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefinition(t *testing.T) {
	t.Parallel()

	tests := []struct {
		scenario string
		function func(*testing.T)
	}{
		{
			scenario: "new definitions",
			function: testNewDefinitions,
		},
		{
			scenario: "successful validation",
			function: testSuccessfulValidation,
		},
		{
			scenario: "empty listen path validation",
			function: testEmptyListenPathValidation,
		},
		{
			scenario: "invalid target url validation",
			function: testInvalidTargetURLValidation,
		},
		{
			scenario: "is balancer defined",
			function: testIsBalancerDefined,
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			test.function(t)
		})
	}
}

func testNewDefinitions(t *testing.T) {
	definition := NewDefinition()

	assert.Equal(t, []string{"GET"}, definition.Methods)
	assert.NotNil(t, definition)
}

func testSuccessfulValidation(t *testing.T) {
	definition := Definition{
		ListenPath: "/*",
		Upstreams: &Upstreams{
			Balancing: "roundrobin",
			Targets: Targets{
				{Target: "http://test.com"},
			},
		},
	}
	isValid, err := definition.Validate()

	assert.NoError(t, err)
	assert.True(t, isValid)
}

func testEmptyListenPathValidation(t *testing.T) {
	definition := Definition{}
	isValid, err := definition.Validate()

	assert.Error(t, err)
	assert.False(t, isValid)
}

func testInvalidTargetURLValidation(t *testing.T) {
	definition := Definition{
		ListenPath: " ",
		Upstreams: &Upstreams{
			Balancing: "roundrobin",
			Targets: Targets{
				{Target: "wrong"},
			},
		},
	}
	isValid, err := definition.Validate()

	assert.Error(t, err)
	assert.False(t, isValid)
}

func testIsBalancerDefined(t *testing.T) {
	definition := NewDefinition()
	assert.False(t, definition.IsBalancerDefined())

	target := &Target{Target: "http://localhost:8080/api-name"}
	definition.Upstreams.Targets = append(definition.Upstreams.Targets, target)
	assert.True(t, definition.IsBalancerDefined())
}
