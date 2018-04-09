package loader

import (
	"testing"

	"net/http"

	"github.com/hellofresh/janus/pkg/api"
	"github.com/hellofresh/janus/pkg/proxy"
	"github.com/hellofresh/janus/pkg/router"
	"github.com/stretchr/testify/assert"
)

func TestLoadAPIsWithParams(t *testing.T) {
	r := router.NewChiRouter()
	Load(proxy.NewRegister(r, proxy.Params{}), api.NewInMemoryRepository())

	assert.Equal(t, 1, r.RoutesCount())
}

func TestLoadValidAPIDefinitions(t *testing.T) {
	r := router.NewChiRouter()

	apiRepo := api.NewInMemoryRepository()
	apiRepo.Add(&api.Definition{
		Name:   "test1",
		Active: true,
		Proxy: &proxy.Definition{
			ListenPath: "/test1",
			Upstreams: &proxy.Upstreams{
				Balancing: "roundrobin",
				Targets: []*proxy.Target{
					&proxy.Target{Target: "http://test1"},
				},
			},
			Methods: []string{http.MethodGet},
		},
		Plugins: []api.Plugin{
			{
				Name:    "oauth2",
				Enabled: false,
			},
			{
				Name:    "compression",
				Enabled: true,
			},
			{
				Name:    "rate_limit",
				Enabled: true,
				Config:  map[string]interface{}{"limit": "10-S", "policy": "local"},
			},
		},
	})
	apiRepo.Add(&api.Definition{
		Name:   "test2",
		Active: true,
		Proxy: &proxy.Definition{
			ListenPath: "/test2",
			Upstreams: &proxy.Upstreams{
				Balancing: "roundrobin",
				Targets: []*proxy.Target{
					&proxy.Target{Target: "http://test2"},
				},
			},
			Methods: []string{http.MethodGet},
		},
	})

	Load(proxy.NewRegister(r, proxy.Params{}), apiRepo)

	assert.Equal(t, 2, r.RoutesCount())
}

func TestLoadInvalidAPIDefinitions(t *testing.T) {
	r := router.NewChiRouter()

	apiRepo := api.NewInMemoryRepository()
	definition := &api.Definition{
		Name:   "test2",
		Active: true,
		Proxy: &proxy.Definition{
			ListenPath: "/test2",
			Upstreams: &proxy.Upstreams{
				Balancing: "roundrobin",
				Targets: []*proxy.Target{
					&proxy.Target{Target: "http://test2"},
				},
			},
			Methods: []string{http.MethodGet},
		},
	}
	err := apiRepo.Add(definition)
	assert.NoError(t, err)

	definition.Name = ""
	Load(proxy.NewRegister(r, proxy.Params{}), apiRepo)

	assert.Equal(t, 1, r.RoutesCount())
}

func TestLoadAPIDefinitionsMissingHTTPMethods(t *testing.T) {
	r := router.NewChiRouter()

	apiRepo := api.NewInMemoryRepository()
	apiRepo.Add(&api.Definition{
		Name:   "test1",
		Active: true,
		Proxy: &proxy.Definition{
			ListenPath: "/test1",
			Upstreams: &proxy.Upstreams{
				Balancing: "roundrobin",
				Targets: []*proxy.Target{
					&proxy.Target{Target: "http://test1"},
				},
			},
		},
	})

	Load(proxy.NewRegister(r, proxy.Params{}), apiRepo)

	assert.Equal(t, 1, r.RoutesCount())
}

func TestLoadInactiveAPIDefinitions(t *testing.T) {
	r := router.NewChiRouter()

	apiRepo := api.NewInMemoryRepository()
	apiRepo.Add(&api.Definition{
		Name:   "test1",
		Active: false,
		Proxy: &proxy.Definition{
			ListenPath: "/test1",
			Upstreams: &proxy.Upstreams{
				Balancing: "roundrobin",
				Targets: []*proxy.Target{
					&proxy.Target{Target: "http://test1"},
				},
			},
		},
	})

	Load(proxy.NewRegister(r, proxy.Params{}), apiRepo)

	assert.Equal(t, 1, r.RoutesCount())
}
