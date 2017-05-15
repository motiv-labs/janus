package proxy

import (
	"io/ioutil"
	"net/http"
	"testing"

	log "github.com/Sirupsen/logrus"
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
	log.SetOutput(ioutil.Discard)

	ts := test.NewServer(createRegisterAndRouter())
	defer ts.Close()

	for _, tc := range tests {
		res, err := ts.Do(tc.method, tc.url, make(map[string]string))
		assert.NoError(t, err)
		if res != nil {
			defer res.Body.Close()
		}

		assert.Equal(t, tc.expectedContentType, res.Header.Get("Content-Type"))
		assert.Equal(t, tc.expectedCode, res.StatusCode, tc.description)
	}
}

func createProxyDefinitions() []*Definition {
	return []*Definition{
		&Definition{
			ListenPath:  "/example/*",
			UpstreamURL: "http://www.mocky.io/v2/58c6c60710000040151b7cad",
			Methods:     []string{"ALL"},
		},
		&Definition{
			ListenPath:  "/posts/*",
			UpstreamURL: "https://jsonplaceholder.typicode.com/posts",
			StripPath:   true,
			Methods:     []string{"ALL"},
		},
		&Definition{
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
	return router.NewChiRouter()
}

func createRegister(r router.Router) *Register {
	var routes []*Route

	definitions := createProxyDefinitions()
	for _, def := range definitions {
		routes = append(routes, NewRoute(def))
	}

	register := NewRegister(r, createProxy())
	register.AddMany(routes)

	return register
}

func createProxy() *Proxy {
	return WithParams(Params{
		StatsClient: stats.NewStatsdClient("", ""),
	})
}
