package main

import (
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/etcinit/speedbump"
	"github.com/gin-gonic/gin"
	"github.com/hellofresh/ginger-middleware/mongodb"
	"gopkg.in/redis.v3"
)

// APIManager handles the creation and configuration of an api definition
type APIManager struct {
	proxyRegister *ProxyRegister
	redisClient   *redis.Client
	accessor      *mongodb.DatabaseAccessor
}

// NewAPIManager creates a new instance of the api manager
func NewAPIManager(router *gin.Engine, redisClient *redis.Client, accessor *mongodb.DatabaseAccessor) *APIManager {
	proxyRegister := &ProxyRegister{Engine: router}
	return &APIManager{proxyRegister, redisClient, accessor}
}

// Load loads all api specs from a datasource
func (m *APIManager) Load() {
	specs := m.getAPISpecs()
	go m.LoadApps(specs)
}

// LoadApps load application middleware
func (m *APIManager) LoadApps(apiSpecs []*APISpec) {
	log.Debug("Loading API configurations")

	for _, referenceSpec := range apiSpecs {
		var skip bool

		//Validates the proxy
		skip = m.validateProxy(referenceSpec.Proxy)

		if false == referenceSpec.Active {
			log.Debug("API is not active, skiping...")
			skip = false
		}

		if skip {
			hasher := speedbump.PerSecondHasher{}
			limit := referenceSpec.RateLimit.Limit
			limiter := speedbump.NewLimiter(m.redisClient, hasher, limit)

			mw := &Middleware{referenceSpec}
			var beforeHandlers = []gin.HandlerFunc{
				CreateMiddleware(&RateLimitMiddleware{mw, limiter, hasher, limit}),
				CreateMiddleware(&CorsMiddleware{mw}),
			}

			if referenceSpec.UseOauth2 {
				log.Debug("Loading OAuth Manager")
				referenceSpec.OAuthManager = &OAuthManager{m.redisClient}
				m.addOAuthHandlers(mw)
				beforeHandlers = append(beforeHandlers, CreateMiddleware(&Oauth2KeyExists{mw}))
				log.Debug("Done loading OAuth Manager")
			}

			m.proxyRegister.Register(referenceSpec.Proxy, beforeHandlers, nil)
			log.Debug("Proxy registered")
		} else {
			log.Error("Listen path is empty, skipping...")
		}
	}
}

//addOAuthHandlers loads configured oauth endpoints
func (m *APIManager) addOAuthHandlers(mw *Middleware) {
	log.Debug("Loading oauth configuration")
	var proxies []Proxy
	var handlers []gin.HandlerFunc

	oauthMeta := mw.Spec.Oauth2Meta

	//oauth proxy
	log.Debug("Registering authorize endpoint")
	authorizeProxy := oauthMeta.OauthEndpoints.Authorize
	if m.validateProxy(authorizeProxy) {
		proxies = append(proxies, authorizeProxy)
	} else {
		log.Debug("No authorize endpoint")
	}

	log.Debug("Registering token endpoint")
	tokenProxy := oauthMeta.OauthEndpoints.Token
	if m.validateProxy(tokenProxy) {
		proxies = append(proxies, tokenProxy)
	} else {
		log.Debug("No token endpoint")
	}

	log.Debug("Registering info endpoint")
	infoProxy := oauthMeta.OauthEndpoints.Info
	if m.validateProxy(infoProxy) {
		proxies = append(proxies, infoProxy)
	} else {
		log.Debug("No info endpoint")
	}

	log.Debug("Registering create client endpoint")
	createProxy := oauthMeta.OauthClientEndpoints.Create
	if m.validateProxy(createProxy) {
		proxies = append(proxies, createProxy)
	} else {
		log.Debug("No client create endpoint")
	}

	log.Debug("Registering remove client endpoint")
	removeProxy := oauthMeta.OauthClientEndpoints.Remove
	if m.validateProxy(removeProxy) {
		proxies = append(proxies, removeProxy)
	} else {
		log.Debug("No client remove endpoint")
	}

	handlers = append(handlers, CreateMiddleware(&OAuthMiddleware{mw}))
	m.proxyRegister.RegisterMany(proxies, nil, handlers)
}

//getAPISpecs Load application specs from datasource
func (m *APIManager) getAPISpecs() []*APISpec {
	log.Debug("Using App Configuration from Mongo DB")
	return APILoader.LoadDefinitionsFromDatastore(m.accessor.Session)
}

//validateProxy validates proxy data
func (m *APIManager) validateProxy(proxy Proxy) bool {
	if proxy.ListenPath == "" {
		log.Error("Listen path is empty")
		return false
	}

	if strings.Contains(proxy.ListenPath, " ") {
		log.Error("Listen path contains spaces, is invalid")
		return false
	}

	return true
}
