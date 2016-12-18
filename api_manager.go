package janus

import (
	log "github.com/Sirupsen/logrus"
	"github.com/etcinit/speedbump"
	"github.com/hellofresh/janus/middleware"
	"github.com/hellofresh/janus/oauth"
	"github.com/hellofresh/janus/router"
	"gopkg.in/redis.v3"
)

var APILoader = APIDefinitionLoader{}

type APIManager struct {
	proxyRegister *ProxyRegister
	redisClient   *redis.Client
	accessor      *middleware.DatabaseAccessor
}

// NewAPIManager creates a new instance of the api manager
func NewAPIManager(router router.Router, redisClient *redis.Client, accessor *middleware.DatabaseAccessor, proxyRegister *ProxyRegister) *APIManager {
	return &APIManager{proxyRegister, redisClient, accessor}
}

// Load loads all api specs from a datasource
func (m *APIManager) Load() {
	oauthManager := &oauth.OAuthManager{m.redisClient}

	oAuthServers := m.getOAuthServers()
	go m.LoadOAuthServers(oAuthServers, oauthManager)

	specs := m.getAPISpecs()
	go m.LoadApps(specs, oauthManager)
}

// LoadApps load application middleware
func (m *APIManager) LoadApps(apiSpecs []*APISpec, oauthManager *oauth.OAuthManager) {
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

			var handlers []router.Constructor
			if referenceSpec.RateLimit.Enabled {
				handlers = append(handlers, NewRateLimitMiddleware(limiter, hasher, limit).Handler)
			} else {
				log.Debug("Rate limit is not enabled")
			}

			if referenceSpec.CorsMeta.Enabled {
				handlers = append(handlers, NewCorsMiddleware(referenceSpec.CorsMeta).Handler)
			} else {
				log.Debug("CORS is not enabled")
			}

			if referenceSpec.UseOauth2 {
				handlers = append(handlers, NewOauth2KeyExistsMiddleware(oauthManager).Handler)
			} else {
				log.Debug("OAuth2 is not enabled")
			}

			m.proxyRegister.Register(referenceSpec.Proxy, handlers...)

			log.Debug("Proxy registered")
		} else {
			log.Error("Listen path is empty, skipping...")
		}
	}
}

// LoadOAuthServers loads and register the oauth servers
func (m *APIManager) LoadOAuthServers(oauthServers []*OAuthSpec, oauthManager *oauth.OAuthManager) {
	log.Debug("Loading OAuth servers configurations")
	oauthRegister := &OAuthRegister{}

	for _, oauthServer := range oauthServers {
		proxies := oauthRegister.GetProxiesForServer(oauthServer.OAuth)
		m.proxyRegister.RegisterMany(
			proxies,
			NewOauth2SecretMiddleware(oauthServer).Handler,
			NewCorsMiddleware(oauthServer.CorsMeta).Handler,
		)
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
