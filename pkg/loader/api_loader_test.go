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
	"github.com/hellofresh/stats-go"
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
	cbFunc          func(t *testing.T) ConfigureCircuitBreakerFunc
}{
	{
		description: "Get example route",
		method:      "GET",
		url:         "/example",
		expectedHeaders: map[string]string{
			"Content-Type": "application/json; charset=utf-8",
		},
		expectedCode: http.StatusOK,
		cbFunc: func(t *testing.T) ConfigureCircuitBreakerFunc {
			return func(name string, c CircuitBreakerConfig) {
				assert.NotEmpty(t, name, "circuit breaker not should not be empty")
				assert.Equal(t, 1, c.Timeout, "unexpected circuit breaker timeout")
				assert.Equal(t, 2, c.MaxConcurrentRequests, "unexpected circuit breaker max concurrent requests")
				assert.Equal(t, 3, c.RequestVolumeThreshold, "unexpected circuit breaker request volume threshold")
				assert.Equal(t, 4, c.SleepWindow, "unexpected circuit breaker sleep window")
				assert.Equal(t, 5, c.ErrorPercentThreshold, "unexpected circuit breaker error percent threshold")
			}
		},
	}, {
		description: "Get invalid route",
		method:      "GET",
		url:         "/invalid-route",
		expectedHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		expectedCode: http.StatusNotFound,
		cbFunc: func(t *testing.T) ConfigureCircuitBreakerFunc {
			return func(name string, c CircuitBreakerConfig) {}
		},
	},
}

func TestSuccessfulLoader(t *testing.T) {
	log.SetOutput(ioutil.Discard)

	for _, tc := range tests {
		routerInstance, err := createRegisterAndRouter(tc.cbFunc(t))
		assert.NoError(t, err)
		ts := test.NewServer(routerInstance)
		defer ts.Close()

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

func createRegisterAndRouter(fn ConfigureCircuitBreakerFunc) (router.Router, error) {
	r := createRouter()
	r.Use(middleware.NewRecovery(errors.RecoveryHandler))

	statsClient, _ := stats.NewClient("noop://", "")
	register := proxy.NewRegister(r, proxy.Params{StatsClient: statsClient})
	proxyRepo, err := api.NewFileSystemRepository("../../assets/apis")
	if err != nil {
		return nil, err
	}

	loader := NewAPILoader(register, ConfigureCircuitBreaker(fn))
	loader.LoadDefinitions(proxyRepo)

	return r, nil
}

func createRouter() router.Router {
	router.DefaultOptions.NotFoundHandler = errors.NotFound
	return router.NewChiRouterWithOptions(router.DefaultOptions)
}
