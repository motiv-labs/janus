// +build integration

package loader

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/hellofresh/janus/pkg/api"
	"github.com/hellofresh/janus/pkg/middleware"
	"github.com/hellofresh/janus/pkg/plugin"
	"github.com/hellofresh/janus/pkg/proxy"
	"github.com/hellofresh/janus/pkg/router"
	"github.com/hellofresh/janus/pkg/test"
	"github.com/hellofresh/janus/pkg/web"
	stats "github.com/hellofresh/stats-go"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var tests = []struct {
	description     string
	method          string
	url             string
	headers         map[string]string
	expectedHeaders map[string]string
	expectedCode    int
}{
	{
		description: "Get example route",
		method:      "GET",
		url:         "/example",
		expectedHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		expectedCode: http.StatusOK,
	}, {
		description: "Get invalid route",
		method:      "GET",
		url:         "/invalid-route",
		expectedHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		expectedCode: http.StatusNotFound,
	},
}

func TestSuccessfulLoader(t *testing.T) {
	log.SetOutput(ioutil.Discard)

	router, err := createRegisterAndRouter()
	assert.NoError(t, err)
	ts := test.NewServer(router)
	defer ts.Close()

	for _, tc := range tests {
		res, err := ts.Do(tc.method, tc.url, tc.headers)
		assert.NoError(t, err)
		if res != nil {
			defer res.Body.Close()
		}

		for headerName, headerValue := range tc.expectedHeaders {
			assert.Equal(t, headerValue, res.Header.Get(headerName))
		}

		assert.Equal(t, tc.expectedCode, res.StatusCode, tc.description)
	}
}

func createRegisterAndRouter() (router.Router, error) {
	r := createRouter()
	r.Use(middleware.NewRecovery(web.RecoveryHandler).Handler)

	register := proxy.NewRegister(r, createProxy())
	proxyRepo, err := createProxyRepo()
	if err != nil {
		return nil, err
	}

	pluginLoader := plugin.NewLoader()
	loader := NewAPILoader(register, pluginLoader)
	loader.LoadDefinitions(proxyRepo)

	return r, nil
}

func createProxyRepo() (api.Repository, error) {
	return api.NewFileSystemRepository("../../examples/apis")
}

func createRouter() router.Router {
	router.DefaultOptions.NotFoundHandler = web.NotFound
	return router.NewChiRouterWithOptions(router.DefaultOptions)
}

func createProxy() *proxy.Proxy {
	return proxy.WithParams(proxy.Params{
		StatsClient: stats.NewStatsdClient("", ""),
	})
}
