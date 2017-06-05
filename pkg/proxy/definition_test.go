package proxy_test

import (
	"testing"

	"github.com/hellofresh/janus/pkg/proxy"
	"github.com/stretchr/testify/assert"
)

func TestSuccessfulValidation(t *testing.T) {
	definition := proxy.Definition{
		ListenPath:  "/*",
		UpstreamURL: "http://test.com",
	}
	isValid, err := definition.Validate()

	assert.NoError(t, err)
	assert.True(t, isValid)
}

func TestEmptyListenPathValidation(t *testing.T) {
	definition := proxy.Definition{}
	isValid, err := definition.Validate()

	assert.Error(t, err)
	assert.False(t, isValid)
}

func TestInvalidTargetURLValidation(t *testing.T) {
	definition := proxy.Definition{
		ListenPath:  " ",
		UpstreamURL: "wrong",
	}
	isValid, err := definition.Validate()

	assert.Error(t, err)
	assert.False(t, isValid)
}

func TestRouteToJSON(t *testing.T) {
	definition := proxy.Definition{
		Methods: make([]string, 0),
		Hosts:   make([]string, 0),
	}
	route := proxy.NewRoute(&definition)
	json, err := route.JSONMarshal()
	assert.NoError(t, err)
	assert.JSONEq(
		t,
		`{"proxy": {"insecure_skip_verify": false, "append_path":false, "enable_load_balancing":false, "methods":[], "hosts":[], "preserve_host":false, "listen_path":"", "upstream_url":"", "strip_path":false}}`,
		string(json),
	)
}

func TestJSONToRoute(t *testing.T) {
	route, err := proxy.JSONUnmarshalRoute([]byte(`{"proxy": {"insecure_skip_verify": false, "append_path":false, "enable_load_balancing":false, "methods":[], "hosts":[], "preserve_host":false, "listen_path":"", "upstream_url":"/*", "strip_path":false}}`))

	assert.NoError(t, err)
	assert.IsType(t, &proxy.Route{}, route)
}

func TestJSONToRouteError(t *testing.T) {
	_, err := proxy.JSONUnmarshalRoute([]byte{})

	assert.Error(t, err)
}
