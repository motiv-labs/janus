package loader

import (
	"github.com/hellofresh/janus/pkg/api"
	"github.com/hellofresh/janus/pkg/oauth"
	"github.com/hellofresh/janus/pkg/plugin"
	"github.com/hellofresh/janus/pkg/proxy"
	"github.com/hellofresh/janus/pkg/router"
	"github.com/hellofresh/janus/pkg/store"
	"github.com/hellofresh/janus/pkg/web"
	stats "github.com/hellofresh/stats-go"

	// this is needed to call the init function on each plugin
	_ "github.com/hellofresh/janus/pkg/plugin/compression"
	_ "github.com/hellofresh/janus/pkg/plugin/cors"
	_ "github.com/hellofresh/janus/pkg/plugin/oauth2"
	_ "github.com/hellofresh/janus/pkg/plugin/rate"
	_ "github.com/hellofresh/janus/pkg/plugin/requesttransformer"
)

// Params initialization options.
type Params struct {
	Router      router.Router
	Storage     store.Store
	APIRepo     api.Repository
	OAuthRepo   oauth.Repository
	StatsClient stats.Client
	ProxyParams proxy.Params
}

// Load loads all the basic components and definitions into a router
func Load(params Params) {
	// create proxy register
	register := proxy.NewRegister(params.Router, params.ProxyParams)

	apiLoader := NewAPILoader(register, plugin.Params{
		Router:      params.Router,
		Storage:     params.Storage,
		APIRepo:     params.APIRepo,
		OAuthRepo:   params.OAuthRepo,
		StatsClient: params.StatsClient,
	})
	apiLoader.LoadDefinitions(params.APIRepo)

	oauthLoader := NewOAuthLoader(register, params.Storage)
	oauthLoader.LoadDefinitions(params.OAuthRepo)

	// some routers may panic when have empty routes list, so add one dummy 404 route to avoid this
	if params.Router.RoutesCount() < 1 {
		params.Router.Any("/", web.NotFound)
	}
}
