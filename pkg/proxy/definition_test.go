package proxy

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/hellofresh/janus/pkg/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		{
			scenario: "add middleware",
			function: testAddMiddlewares,
		},
		{
			scenario: "unmarshal forwarding_timeouts from json",
			function: testUnmarshalForwardingTimeoutsFromJSON,
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
	assert.Len(t, definition.Upstreams.Targets.ToBalancerTargets(), 1)
}

func testAddMiddlewares(t *testing.T) {
	routerDefinition := NewRouterDefinition(NewDefinition())
	routerDefinition.AddMiddleware(middleware.NewLogger().Handler)

	assert.Len(t, routerDefinition.Middleware(), 1)
}

func testUnmarshalForwardingTimeoutsFromJSON(t *testing.T) {
	rawDefinition := []byte(`
  {
    "preserve_host":false,
    "listen_path":"/example/*",
    "upstreams":{
      "balancing":"roundrobin",
      "targets":[
        {
          "target":"http://localhost:9089/hello-world"
        }
      ]
    },
    "strip_path":false,
    "append_path":false,
    "methods":[
      "GET"
    ],
    "forwarding_timeouts": {
      "dial_timeout": "30s",
      "response_header_timeout": "31s"
    }
  }
`)
	definition := NewDefinition()
	err := json.Unmarshal(rawDefinition, &definition)
	require.NoError(t, err)

	assert.Equal(t, 30*time.Second, time.Duration(definition.ForwardingTimeouts.DialTimeout))
	assert.Equal(t, 31*time.Second, time.Duration(definition.ForwardingTimeouts.ResponseHeaderTimeout))
}
