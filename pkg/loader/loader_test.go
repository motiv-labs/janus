package loader

import (
	"testing"

	"net/http"

	"github.com/hellofresh/janus/pkg/api"
	"github.com/hellofresh/janus/pkg/oauth"
	"github.com/hellofresh/janus/pkg/proxy"
	"github.com/hellofresh/janus/pkg/router"
	"github.com/hellofresh/janus/pkg/store"
	stats "github.com/hellofresh/stats-go"
	"github.com/stretchr/testify/assert"
)

func TestLoadAPIsWithParams(t *testing.T) {
	r := router.NewChiRouter()
	Load(loadParamsForTest(r, api.NewInMemoryRepository()))

	assert.Equal(t, 1, r.RoutesCount())
}

func TestLoadValidAPIDefinitions(t *testing.T) {
	r := router.NewChiRouter()

	apiRepo := api.NewInMemoryRepository()
	apiRepo.Add(&api.Definition{
		Name:   "test1",
		Active: true,
		Proxy: &proxy.Definition{
			ListenPath:  "/test1",
			UpstreamURL: "http://test1",
			Methods:     []string{http.MethodGet},
		},
	})
	apiRepo.Add(&api.Definition{
		Name:   "test2",
		Active: true,
		Proxy: &proxy.Definition{
			ListenPath:  "/test2",
			UpstreamURL: "http://test2",
			Methods:     []string{http.MethodGet},
		},
	})

	Load(loadParamsForTest(r, apiRepo))

	assert.Equal(t, 2, r.RoutesCount())
}

func TestLoadAPIDefinitionsMissingHTTPMethods(t *testing.T) {
	r := router.NewChiRouter()

	apiRepo := api.NewInMemoryRepository()
	apiRepo.Add(&api.Definition{
		Name:   "test1",
		Active: true,
		Proxy: &proxy.Definition{
			ListenPath:  "/test1",
			UpstreamURL: "http://test1",
		},
	})

	Load(loadParamsForTest(r, apiRepo))

	assert.Equal(t, 1, r.RoutesCount())
}

func TestLoadInactiveAPIDefinitions(t *testing.T) {
	r := router.NewChiRouter()

	apiRepo := api.NewInMemoryRepository()
	apiRepo.Add(&api.Definition{
		Name:   "test1",
		Active: false,
		Proxy: &proxy.Definition{
			ListenPath:  "/test1",
			UpstreamURL: "http://test1",
		},
	})

	Load(loadParamsForTest(r, apiRepo))

	assert.Equal(t, 1, r.RoutesCount())
}

func loadParamsForTest(r router.Router, apiRepo api.Repository) Params {
	return Params{
		Storage:   store.NewInMemoryStore(),
		APIRepo:   apiRepo,
		OAuthRepo: oauth.NewInMemoryRepository(),
		Router:    r,
		ProxyParams: proxy.Params{
			StatsClient: stats.NewStatsdClient("", ""),
		},
	}

}
