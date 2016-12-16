package janus

import (
	log "github.com/Sirupsen/logrus"
	"github.com/etcinit/speedbump"
	"github.com/hellofresh/ginger-middleware/mongodb"
	"gopkg.in/alexcesaro/statsd.v2"
	"gopkg.in/redis.v3"
)

var APILoader = APIDefinitionLoader{}

type APIManager struct {
	proxyRegister *ProxyRegister
	redisClient   *redis.Client
	accessor      *mongodb.DatabaseAccessor
}

// NewAPIManager creates a new instance of the api manager
func NewAPIManager(router Router, redisClient *redis.Client, accessor *mongodb.DatabaseAccessor, statsdClient *statsd.Client) *APIManager {
	proxyRegister := &ProxyRegister{Router: router, statsdClient: statsdClient}
	return &APIManager{proxyRegister, redisClient, accessor}
}

// Load loads all api specs from a datasource
func (m *APIManager) Load() {
	oauthManager := &OAuthManager{m.redisClient}

	oAuthServers := m.getOAuthServers()
	go m.LoadOAuthServers(oAuthServers, oauthManager)

	specs := m.getAPISpecs()
	go m.LoadApps(specs, oauthManager)
}

// LoadApps load application middleware
func (m *APIManager) LoadApps(apiSpecs []*APISpec, oauthManager *OAuthManager) {
	log.Debug("Loading API configurations")

	for _, referenceSpec := range apiSpecs {
		var skip bool

		//Validates the proxy
		skip = validateProxy(referenceSpec.Proxy)
		if false == referenceSpec.Active {
			log.Debug("API is not active, skiping...")
			skip = false
		}

		if skip {
			hasher := speedbump.PerSecondHasher{}
			limit := referenceSpec.RateLimit.Limit
			limiter := speedbump.NewLimiter(m.redisClient, hasher, limit)

			mw := &Middleware{referenceSpec}
			var beforeHandlers = []HandlerFunc{
				CreateMiddleware(&RateLimitMiddleware{mw, limiter, hasher, limit}),
				CreateMiddleware(&CorsMiddleware{mw}),
			}

			if referenceSpec.UseOauth2 {
				beforeHandlers = append(beforeHandlers, CreateMiddleware(&Oauth2KeyExistsMiddleware{mw, oauthManager}))
			}

			m.proxyRegister.Register(referenceSpec.Proxy, beforeHandlers, nil)
			log.Debug("Proxy registered")
		} else {
			log.Error("Listen path is empty, skipping...")
		}
	}
}

// LoadOAuthServers loads and register the oauth servers
func (m *APIManager) LoadOAuthServers(oauthServers []*OAuthSpec, oauthManager *OAuthManager) {
	log.Debug("Loading OAuth servers configurations")

	var beforeHandlers []HandlerFunc
	var handlers []HandlerFunc
	oauthRegister := &OAuthRegister{}

	for _, oauthServer := range oauthServers {
		beforeHandlers = append(beforeHandlers, CreateMiddleware(&Oauth2SecretMiddleware{oauthServer}))
		handlers = append(handlers, CreateMiddleware(&OAuthMiddleware{oauthManager, oauthServer}))
		oauthServer.OAuthManager = &OAuthManager{m.redisClient}
		proxies := oauthRegister.GetProxiesForServer(oauthServer.OAuth)
		m.proxyRegister.RegisterMany(proxies, beforeHandlers, handlers)
	}

	log.Debug("Done loading OAuth servers configurations")
}

//getAPISpecs Load application specs from datasource
func (m *APIManager) getAPISpecs() []*APISpec {
	log.Debug("Using App Configuration from Mongo DB")
	return APILoader.LoadDefinitionsFromDatastore(m.accessor.Session)
}

//getOAuthServers Load oauth servers from datasource
func (m *APIManager) getOAuthServers() []*OAuthSpec {
	log.Debug("Using Oauth servers configuration from Mongo DB")
	return APILoader.LoadOauthServersFromDatastore(m.accessor.Session)
}
