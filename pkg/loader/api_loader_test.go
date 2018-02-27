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
				if name == "" {
					t.Fatalf("unexpected circuit breaker timeout: %d", c.Timeout)
				}
				if c.Timeout != 1 {
					t.Fatalf("unexpected circuit breaker timeout: %d", c.Timeout)
				}
				if c.MaxConcurrentRequests != 2 {
					t.Fatalf("unexpected circuit breaker max concurrent requests: %d", c.MaxConcurrentRequests)
				}
				if c.RequestVolumeThreshold != 3 {
					t.Fatalf("unexpected circuit breaker request volume threshold: %d", c.RequestVolumeThreshold)
				}
				if c.SleepWindow != 4 {
					t.Fatalf("unexpected circuit breaker sleep window: %d", c.SleepWindow)
				}
				if c.ErrorPercentThreshold != 5 {
					t.Fatalf("unexpected circuit breaker error percent threshold: %d", c.ErrorPercentThreshold)
				}
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
