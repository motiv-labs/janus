// +build integration

package proxy

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/hellofresh/stats-go/client"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hellofresh/janus/pkg/router"
	"github.com/hellofresh/janus/pkg/test"
)

const defaultUpstreamsPort = "9089"

var tests = []struct {
	description         string
	method              string
	url                 string
	expectedContentType string
	expectedCode        int
}{
	{
		description:         "Get example route",
		method:              "GET",
		url:                 "/example",
		expectedContentType: "application/json; charset=utf-8",
		expectedCode:        http.StatusOK,
	}, {
		description:         "Get invalid route",
		method:              "GET",
		url:                 "/invalid-route",
		expectedContentType: "text/plain; charset=utf-8",
		expectedCode:        http.StatusNotFound,
	},
	{
		description:         "Get one posts - strip path",
		method:              "GET",
		url:                 "/posts/1",
		expectedContentType: "application/json; charset=utf-8",
		expectedCode:        http.StatusOK,
	},
	{
		description:         "Get one posts - append path",
		method:              "GET",
		url:                 "/append",
		expectedContentType: "application/json; charset=utf-8",
		expectedCode:        http.StatusOK,
	},
	{
		description:         "Get one recipe - parameter interpolation",
		url:                 "/api/recipes/5252b1b5301bbf46038b473f",
		expectedContentType: "application/json; charset=utf-8",
		expectedCode:        http.StatusOK,
	},
	{
		description:         "No parameter to interpolate",
		url:                 "/api/recipes/search",
		expectedContentType: "application/json",
		expectedCode:        http.StatusNotFound,
	},
	{
		description:         "named param with strip path to posts",
		method:              "GET",
		url:                 "/private/localhost/posts/1",
		expectedContentType: "application/json; charset=utf-8",
		expectedCode:        http.StatusOK,
	},
	{
		description:         "named param with strip path - no api",
		method:              "GET",
		url:                 "/private/localhost/no-api",
		expectedContentType: "application/json",
		expectedCode:        http.StatusNotFound,
	},
}

func TestSuccessfulProxy(t *testing.T) {
	t.Parallel()

	log.SetOutput(ioutil.Discard)

	ts := test.NewServer(createRegisterAndRouter(t))
	defer ts.Close()

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			res, err := ts.Do(tc.method, tc.url, make(map[string]string))
			require.NoError(t, err)
			t.Cleanup(func() {
				err := res.Body.Close()
				assert.NoError(t, err)
			})

			assert.Equal(t, tc.expectedContentType, res.Header.Get("Content-Type"), tc.description)
			assert.Equal(t, tc.expectedCode, res.StatusCode, tc.description)
		})
	}
}

func createProxyDefinitions(t *testing.T) []*Definition {
	t.Helper()

	upstreamsPort := os.Getenv("DYNAMIC_UPSTREAMS_PORT")
	if upstreamsPort == "" {
		upstreamsPort = defaultUpstreamsPort
	}

	return []*Definition{
		{
			ListenPath: "/example/*",
			Upstreams: &Upstreams{
				Balancing: "roundrobin",
				Targets:   []*Target{{Target: fmt.Sprintf("http://localhost:%s/hello-world", upstreamsPort)}},
			},
			Methods: []string{"ALL"},
		},
		{
			ListenPath: "/posts/*",
			StripPath:  true,
			Upstreams: &Upstreams{
				Balancing: "roundrobin",
				Targets:   []*Target{{Target: fmt.Sprintf("http://localhost:%s/posts", upstreamsPort)}},
			},
			Methods: []string{"ALL"},
		},
		{
			ListenPath: "/append/*",
			Upstreams: &Upstreams{
				Balancing: "roundrobin",
				Targets:   []*Target{{Target: fmt.Sprintf("http://localhost:%s/hello-world", upstreamsPort)}},
			},
			AppendPath: true,
			Methods:    []string{"GET"},
		},
		{
			ListenPath: "/api/recipes/{id:[\\da-f]{24}}",
			Upstreams: &Upstreams{
				Balancing: "roundrobin",
				Targets:   []*Target{{Target: fmt.Sprintf("http://localhost:%s/recipes/{id}", upstreamsPort)}},
			},
			Methods: []string{"GET"},
		},
		{
			ListenPath: "/api/recipes/search",
			Upstreams: &Upstreams{
				Balancing: "roundrobin",
				Targets:   []*Target{{Target: fmt.Sprintf("http://localhost:%s/recipes/{id}", upstreamsPort)}},
			},
			Methods: []string{"GET"},
		},
		{
			ListenPath: "/private/{service}/*",
			Upstreams: &Upstreams{
				Balancing: "roundrobin",
				Targets:   []*Target{{Target: fmt.Sprintf("http://{service}:%s/", upstreamsPort)}},
			},
			StripPath: true,
			Methods: []string{"ALL"},
		},
	}
}

func createRegisterAndRouter(t *testing.T) router.Router {
	t.Helper()

	r := router.NewChiRouter()
	createRegister(t, r)
	return r
}

func createRegister(t *testing.T, r router.Router) *Register {
	t.Helper()

	register := NewRegister(WithRouter(r), WithStatsClient(client.NewNoop()))

	definitions := createProxyDefinitions(t)
	for _, def := range definitions {
		err := register.Add(NewRouterDefinition(def))
		require.NoError(t, err)
	}

	return register
}
