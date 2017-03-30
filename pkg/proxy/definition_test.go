package proxy_test

import (
	"testing"

	"github.com/hellofresh/janus/pkg/proxy"
	"github.com/stretchr/testify/assert"
)

func TestSuccessfulValidation(t *testing.T) {
	definition := proxy.Definition{
		ListenPath: "/*",
	}

	assert.True(t, proxy.Validate(&definition))
}

func TestEmptyListenPathValidation(t *testing.T) {
	definition := proxy.Definition{}

	assert.False(t, proxy.Validate(&definition))
}

func TestSpaceInListenPathValidation(t *testing.T) {
	definition := proxy.Definition{
		ListenPath: " ",
	}

	assert.False(t, proxy.Validate(&definition))
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
		`{"proxy": {"append_path":false, "enable_load_balancing":false, "methods":[], "hosts":[], "preserve_host":false, "listen_path":"", "upstream_url":"", "strip_path":false}}`,
		string(json),
	)
}

func TestJSONToRoute(t *testing.T) {
	route, err := proxy.JSONUnmarshalRoute([]byte(`{"proxy": {"append_path":false, "enable_load_balancing":false, "methods":[], "hosts":[], "preserve_host":false, "listen_path":"", "upstream_url":"/*", "strip_path":false}}`))

	assert.NoError(t, err)
	assert.IsType(t, &proxy.Route{}, route)
}
