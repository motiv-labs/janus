package proxy

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDefinitions(t *testing.T) {
	definition := NewDefinition()
	assert.NotNil(t, definition)
}

func TestSuccessfulValidation(t *testing.T) {
	definition := Definition{
		ListenPath:  "/*",
		UpstreamURL: "http://test.com",
	}
	isValid, err := definition.Validate()

	assert.NoError(t, err)
	assert.True(t, isValid)
}

func TestEmptyListenPathValidation(t *testing.T) {
	definition := Definition{}
	isValid, err := definition.Validate()

	assert.Error(t, err)
	assert.False(t, isValid)
}

func TestInvalidTargetURLValidation(t *testing.T) {
	definition := Definition{
		ListenPath:  " ",
		UpstreamURL: "wrong",
	}
	isValid, err := definition.Validate()

	assert.Error(t, err)
	assert.False(t, isValid)
}

func TestRouteToJSON(t *testing.T) {
	definition := NewDefinition()
	route := NewRoute(definition)
	json, err := route.JSONMarshal()
	assert.NoError(t, err)
	assert.JSONEq(
		t,
		`{"proxy": {"insecure_skip_verify": false, "append_path":false, "enable_load_balancing":false, "methods":[], "hosts":[], "preserve_host":false, "listen_path":"", "upstream_url":"", "strip_path":false, "upstreams": {"balancing": "", "targets": [] }}}`,
		string(json),
	)
}

func TestJSONToRoute(t *testing.T) {
	route, err := JSONUnmarshalRoute([]byte(`{"proxy": {"insecure_skip_verify": false, "append_path":false, "enable_load_balancing":false, "methods":[], "hosts":[], "preserve_host":false, "listen_path":"", "upstream_url":"/*", "strip_path":false}}`))

	assert.NoError(t, err)
	assert.IsType(t, &Route{}, route)
}

func TestJSONToRouteError(t *testing.T) {
	_, err := JSONUnmarshalRoute([]byte{})

	assert.Error(t, err)
}
