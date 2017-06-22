package loader

import (
	"github.com/hellofresh/janus/pkg/api"
	"github.com/hellofresh/janus/pkg/oauth"
	"github.com/hellofresh/janus/pkg/plugin"
	"github.com/hellofresh/janus/pkg/proxy"
	"github.com/hellofresh/janus/pkg/router"
	"github.com/hellofresh/janus/pkg/store"
	"github.com/hellofresh/janus/pkg/web"
)

// Params initialization options.
type Params struct {
	Router      router.Router
	Storage     store.Store
	APIRepo     api.Repository
	OAuthRepo   oauth.Repository
	ProxyParams proxy.Params
}

// Load loads all the basic components and definitions into a router
func Load(params Params) {
	pluginLoader := plugin.NewLoader()
	pluginLoader.Add(
		plugin.NewRateLimit(params.Storage, params.ProxyParams.StatsClient),
		plugin.NewCORS(),
		plugin.NewOAuth2(params.OAuthRepo, params.Storage),
		plugin.NewCompression(),
		plugin.NewRequestTransformer(),
	)

	// create proxy register
	register := proxy.NewRegister(params.Router, params.ProxyParams)

	apiLoader := NewAPILoader(register, pluginLoader)
	apiLoader.LoadDefinitions(params.APIRepo)

	oauthLoader := NewOAuthLoader(register, params.Storage)
	oauthLoader.LoadDefinitions(params.OAuthRepo)

	// some routers may panic when have empty routes list, so add one dummy 404 route to avoid this
	if params.Router.RoutesCount() < 1 {
		params.Router.Any("/", web.NotFound)
	}
}
