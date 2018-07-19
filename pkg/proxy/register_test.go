// +build integration

package proxy

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/hellofresh/janus/pkg/router"
	"github.com/hellofresh/janus/pkg/test"
	"github.com/hellofresh/stats-go/client"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

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
}

func TestSuccessfulProxy(t *testing.T) {
	t.Parallel()

	log.SetOutput(ioutil.Discard)

	ts := test.NewServer(createRegisterAndRouter())
	defer ts.Close()

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			res, err := ts.Do(tc.method, tc.url, make(map[string]string))
			assert.NoError(t, err)
			if res != nil {
				defer res.Body.Close()
			}

			assert.Equal(t, tc.expectedContentType, res.Header.Get("Content-Type"), tc.description)
			assert.Equal(t, tc.expectedCode, res.StatusCode, tc.description)
		})
	}
}

func createProxyDefinitions() []*Definition {
	return []*Definition{
		{
			ListenPath: "/example/*",
			Upstreams: &Upstreams{
				Balancing: "roundrobin",
				Targets:   []*Target{{Target: "http://localhost:9089/hello-world"}},
			},
			Methods: []string{"ALL"},
		},
		{
			ListenPath: "/posts/*",
			StripPath:  true,
			Upstreams: &Upstreams{
				Balancing: "roundrobin",
				Targets:   []*Target{{Target: "http://localhost:9089/posts"}},
			},
			Methods: []string{"ALL"},
		},
		{
			ListenPath: "/append/*",
			Upstreams: &Upstreams{
				Balancing: "roundrobin",
				Targets:   []*Target{{Target: "http://localhost:9089/hello-world"}},
			},
			AppendPath: true,
			Methods:    []string{"GET"},
		},
		{
			ListenPath: "/api/recipes/{id:[\\da-f]{24}}",
			Upstreams: &Upstreams{
				Balancing: "roundrobin",
				Targets:   []*Target{{Target: "http://localhost:9089/recipes/{id}"}},
			},
			Methods: []string{"GET"},
		},
		{
			ListenPath: "/api/recipes/search",
			Upstreams: &Upstreams{
				Balancing: "roundrobin",
				Targets:   []*Target{{Target: "http://localhost:9089/recipes/{id}"}},
			},
			Methods: []string{"GET"},
		},
	}
}

func createRegisterAndRouter() router.Router {
	r := router.NewChiRouter()
	createRegister(r)
	return r
}

func createRegister(r router.Router) *Register {
	register := NewRegister(WithRouter(r), WithStatsClient(client.NewNoop()))

	definitions := createProxyDefinitions()
	for _, def := range definitions {
		register.Add(NewRouterDefinition(def))
	}

	return register
}
