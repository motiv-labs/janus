// +build integration

package loader

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/hellofresh/janus/pkg/api"
	"github.com/hellofresh/janus/pkg/errors"
	"github.com/hellofresh/janus/pkg/middleware"
	"github.com/hellofresh/janus/pkg/proxy"
	"github.com/hellofresh/janus/pkg/router"
	"github.com/hellofresh/janus/pkg/test"
	"github.com/hellofresh/stats-go/client"
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
			"Content-Type": "application/json; charset=utf-8",
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

	routerInstance, err := createRegisterAndRouter()
	assert.NoError(t, err)
	ts := test.NewServer(routerInstance)
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
	r.Use(middleware.NewRecovery(errors.RecoveryHandler))

	register := proxy.NewRegister(proxy.WithRouter(r), proxy.WithStatsClient(client.NewNoop(false)))
	proxyRepo, err := api.NewFileSystemRepository("../../assets/apis")
	if err != nil {
		return nil, err
	}
	defs, err := proxyRepo.FindAll()
	if err != nil {
		return nil, err
	}

	loader := NewAPILoader(register)
	loader.RegisterAPIs(defs)

	return r, nil
}

func createRouter() router.Router {
	router.DefaultOptions.NotFoundHandler = errors.NotFound
	return router.NewChiRouterWithOptions(router.DefaultOptions)
}
