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
			scenario: "route to json",
			function: testRouteToJSON,
		},
		{
			scenario: "json to route",
			function: testJSONToRoute,
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
			Targets: []*Target{
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
			Targets: []*Target{
				{Target: "wrong"},
			},
		},
	}
	isValid, err := definition.Validate()

	assert.Error(t, err)
	assert.False(t, isValid)
}

func testRouteToJSON(t *testing.T) {
	expectedJSON := `
	{
		"proxy":{
			"insecure_skip_verify":false,
			"append_path":false,
			"enable_load_balancing":false,
			"methods":[
				"GET"
			],
			"hosts":[

			],
			"preserve_host":false,
			"listen_path":"",
			"strip_path":false,
			"upstreams":{
				"balancing":"",
				"targets":[

				]
			}
		}
	}
	`
	definition := NewDefinition()
	route := NewRoute(definition)
	json, err := route.JSONMarshal()

	assert.NoError(t, err)
	assert.JSONEq(t, expectedJSON, string(json))
}

func testJSONToRoute(t *testing.T) {
	route, err := JSONUnmarshalRoute([]byte(`
	{
		"proxy":{
			"insecure_skip_verify":false,
			"append_path":false,
			"enable_load_balancing":false,
			"methods":[],
			"hosts":[],
			"preserve_host":false,
			"listen_path":"",
			"strip_path":false
		}
	}`))

	assert.NoError(t, err)
	assert.IsType(t, &Route{}, route)
}

func testJSONToRouteError(t *testing.T) {
	_, err := JSONUnmarshalRoute([]byte{})
	assert.Error(t, err)
}

func testIsBalancerDefined(t *testing.T) {
	definition := NewDefinition()
	assert.False(t, definition.IsBalancerDefined())

	target := &Target{Target: "http://localhost:8080/api-name"}
	definition.Upstreams.Targets = append(definition.Upstreams.Targets, target)
	assert.True(t, definition.IsBalancerDefined())
}
