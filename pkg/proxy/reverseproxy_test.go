package proxy_test

import (
	"net/http"
	"testing"

	"github.com/hellofresh/janus/pkg/proxy"
	"github.com/hellofresh/janus/pkg/router"
	"github.com/hellofresh/janus/pkg/test"
	stats "github.com/hellofresh/stats-go"
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
}

func TestSuccessfulProxy(t *testing.T) {
	ts := test.NewServer(createRegisterAndRouter())
	defer ts.Close()

	for _, tc := range tests {
		res, err := ts.Do(tc.method, tc.url)
		assert.NoError(t, err)
		if res != nil {
			defer res.Body.Close()
		}

		// b, err := ioutil.ReadAll(res.Body)
		// assert.NoError(t, err)

		assert.Equal(t, tc.expectedContentType, res.Header.Get("Content-Type"))
		assert.Equal(t, tc.expectedCode, res.StatusCode, tc.description)
	}
}

func createProxyDefinitions() []*proxy.Definition {
	return []*proxy.Definition{
		&proxy.Definition{
			ListenPath:  "/example/*",
			UpstreamURL: "http://www.mocky.io/v2/58c6c60710000040151b7cad",
			Methods:     []string{"ALL"},
		},
		&proxy.Definition{
			ListenPath:  "/posts/*",
			UpstreamURL: "https://jsonplaceholder.typicode.com/posts",
			StripPath:   true,
			Methods:     []string{"GET"},
		},
		&proxy.Definition{
			ListenPath:  "/append/*",
			UpstreamURL: "http://www.mocky.io/v2/58c6c60710000040151b7cad",
			AppendPath:  true,
			Methods:     []string{"GET"},
		},
	}
}

func createRegisterAndRouter() router.Router {
	r := createRouter()
	createRegister(r)
	return r
}

func createRouter() router.Router {
	return router.NewHTTPTreeMuxRouter()
}

func createRegister(r router.Router) *proxy.Register {
	var routes []*proxy.Route

	definitions := createProxyDefinitions()
	for _, def := range definitions {
		routes = append(routes, proxy.NewRoute(def))
	}

	register := proxy.NewRegister(r, createProxy())
	register.AddMany(routes)

	return register
}

func createProxy() *proxy.Proxy {
	return proxy.WithParams(proxy.Params{
		StatsClient: stats.NewStatsdStatsClient("", ""),
	})
}
