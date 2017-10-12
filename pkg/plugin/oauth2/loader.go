package oauth2

import (
	"net/http"

	"github.com/hellofresh/janus/pkg/proxy"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
	"github.com/ulule/limiter"
	"github.com/ulule/limiter/drivers/middleware/stdlib"
	smemory "github.com/ulule/limiter/drivers/store/memory"
)

// OAuthLoader handles the loading of the api specs
type OAuthLoader struct {
	register *proxy.Register
}

// NewOAuthLoader creates a new instance of the Loader
func NewOAuthLoader(register *proxy.Register) *OAuthLoader {
	return &OAuthLoader{register}
}

// LoadDefinitions loads all oauth servers from a data source
func (m *OAuthLoader) LoadDefinitions(repo Repository) {
	oAuthServers := m.getOAuthServers(repo)
	m.RegisterOAuthServers(oAuthServers, repo)
}

// RegisterOAuthServers register many oauth servers
func (m *OAuthLoader) RegisterOAuthServers(oauthServers []*Spec, repo Repository) {
	log.Debug("Loading OAuth servers configurations")

	for _, oauthServer := range oauthServers {
		var corsHandler func(h http.Handler) http.Handler
		var rateLimitHandler func(h http.Handler) http.Handler

		logger := log.WithField("name", oauthServer.Name)
		logger.Debug("Registering OAuth server")

		corsHandler = cors.New(cors.Options{
			AllowedOrigins:   oauthServer.CorsMeta.Domains,
			AllowedMethods:   oauthServer.CorsMeta.Methods,
			AllowedHeaders:   oauthServer.CorsMeta.RequestHeaders,
			ExposedHeaders:   oauthServer.CorsMeta.ExposedHeaders,
			AllowCredentials: true,
		}).Handler

		rate, err := limiter.NewRateFromFormatted(oauthServer.RateLimit.Limit)
		if err != nil {
			logger.WithError(err).Error("Not able to create rate limit")
		}

		limiterStore := smemory.NewStore()
		limiterInstance := limiter.New(limiterStore, rate)
		rateLimitHandler = stdlib.NewMiddleware(limiterInstance).Handler

		endpoints := map[*proxy.Definition]proxy.InChain{
			oauthServer.Endpoints.Authorize:    proxy.NewInChain(corsHandler, rateLimitHandler),
			oauthServer.Endpoints.Token:        proxy.NewInChain(NewSecretMiddleware(oauthServer).Handler, corsHandler, rateLimitHandler),
			oauthServer.Endpoints.Introspect:   proxy.NewInChain(corsHandler, rateLimitHandler),
			oauthServer.Endpoints.Revoke:       proxy.NewInChain(corsHandler, rateLimitHandler),
			oauthServer.ClientEndpoints.Create: proxy.NewInChain(corsHandler, rateLimitHandler),
			oauthServer.ClientEndpoints.Remove: proxy.NewInChain(corsHandler, rateLimitHandler),
		}

		m.registerRoutes(endpoints)
		logger.Debug("OAuth server registered")
	}

	log.Debug("Done loading OAuth servers configurations")
}

func (m *OAuthLoader) getOAuthServers(repo Repository) []*Spec {
	oauthServers, err := repo.FindAll()
	if err != nil {
		log.Panic(err)
	}

	var specs []*Spec
	for _, oauthServer := range oauthServers {
		spec := new(Spec)
		spec.OAuth = oauthServer
		manager, err := m.getManager(oauthServer)
		if nil != err {
			log.WithError(err).Error("Oauth definition is not well configured, skipping...")
			continue
		}
		spec.Manager = manager
		specs = append(specs, spec)
	}

	return specs
}

func (m *OAuthLoader) getManager(oauthServer *OAuth) (Manager, error) {
	managerType, err := ParseType(oauthServer.TokenStrategy.Name)
	if nil != err {
		return nil, err
	}

	return NewManagerFactory(oauthServer).Build(managerType)
}

func (m *OAuthLoader) registerRoutes(endpoints map[*proxy.Definition]proxy.InChain) {
	for endpoint, middleware := range endpoints {
		if endpoint == nil {
			log.Debug("Endpoint not registered")
			continue
		}

		l := log.WithField("listen_path", endpoint.ListenPath)
		l.Debug("Registering OAuth endpoint")
		if isValid, err := endpoint.Validate(); isValid && err == nil {
			m.register.Add(proxy.NewRouteWithInOut(endpoint, middleware, nil))
			l.Debug("Endpoint registered")
		} else {
			l.WithError(err).Error("Error when registering endpoint")
		}
	}
}
