package main

import (
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/etcinit/speedbump"
	"github.com/gin-gonic/gin"
	"gopkg.in/redis.v3"
)

// APIManager handles the creation and configuration of an api definition
type APIManager struct {
	proxyRegister *ProxyRegister
	redisClient   *redis.Client
	accessor      *DatabaseAccessor
}

// NewAPIManager creates a new instance of the api manager
func NewAPIManager(router *gin.Engine, redisClient *redis.Client, accessor *DatabaseAccessor) *APIManager {
	proxyRegister := &ProxyRegister{Engine: router}
	return &APIManager{proxyRegister, redisClient, accessor}
}

// Load loads all api specs from a datasource
func (m APIManager) Load() {
	specs := m.getAPISpecs()
	go m.LoadApps(specs)
}

// LoadApps load application middleware
func (m APIManager) LoadApps(apiSpecs []*APISpec) {
	log.Debug("Loading API configurations")

	for _, referenceSpec := range apiSpecs {
		var skip bool

		//Define fields to log
		logger := createContextualLogger(referenceSpec)

		//Validates the proxy
		skip = m.validateProxy(referenceSpec.Proxy)

		if false == referenceSpec.Active {
			logger.Info("API is not active, skiping...")
			skip = false
		}

		if skip {
			cb := NewCircuitBreaker(referenceSpec)

			hasher := speedbump.PerSecondHasher{}
			limit := referenceSpec.RateLimit.Limit
			limiter := speedbump.NewLimiter(m.redisClient, hasher, limit)

			mw := &Middleware{referenceSpec, logger}
			var beforeHandlers = []gin.HandlerFunc{
				CreateMiddleware(&RateLimitMiddleware{mw, limiter, hasher, limit}),
			}

			if referenceSpec.UseOauth2 {
				logger.Debug("Loading OAuth Manager")
				referenceSpec.OAuthManager = &OAuthManager{m.redisClient}
				m.addOAuthHandlers(mw, cb, logger)
				beforeHandlers = append(beforeHandlers, CreateMiddleware(Oauth2KeyExists{mw}))
				logger.Debug("Done loading OAuth Manager")
			}

			m.proxyRegister.Register(referenceSpec.Proxy, cb, beforeHandlers, nil)
			logger.Debug("Proxy registered")
		} else {
			logger.Error("Listen path is empty, skipping...")
		}
	}
}

//addOAuthHandlers loads configured oauth endpoints
func (m APIManager) addOAuthHandlers(mw *Middleware, cb *ExtendedCircuitBreakerMeta, logger *Logger) {
	logger.Info("Loading oauth configuration")
	var proxies []Proxy
	var handlers []gin.HandlerFunc

	oauthMeta := mw.Spec.Oauth2Meta

	//oauth proxy
	logger.Debug("Registering authorize endpoint")
	authorizeProxy := oauthMeta.OauthEndpoints.Authorize
	if m.validateProxy(authorizeProxy) {
		proxies = append(proxies, authorizeProxy)
	} else {
		logger.Debug("No authorize endpoint")
	}

	logger.Debug("Registering token endpoint")
	tokenProxy := oauthMeta.OauthEndpoints.Token
	if m.validateProxy(tokenProxy) {
		proxies = append(proxies, tokenProxy)
	} else {
		logger.Debug("No token endpoint")
	}

	logger.Debug("Registering info endpoint")
	infoProxy := oauthMeta.OauthEndpoints.Info
	if m.validateProxy(infoProxy) {
		proxies = append(proxies, infoProxy)
	} else {
		logger.Debug("No info endpoint")
	}

	logger.Debug("Registering create client endpoint")
	createProxy := oauthMeta.OauthClientEndpoints.Create
	if m.validateProxy(createProxy) {
		proxies = append(proxies, createProxy)
	} else {
		logger.Debug("No client create endpoint")
	}

	logger.Debug("Registering remove client endpoint")
	removeProxy := oauthMeta.OauthClientEndpoints.Remove
	if m.validateProxy(removeProxy) {
		proxies = append(proxies, removeProxy)
	} else {
		logger.Debug("No client remove endpoint")
	}

	handlers = append(handlers, CreateMiddleware(OAuthMiddleware{mw}))
	m.proxyRegister.registerMany(proxies, cb, nil, handlers)
}

//getAPISpecs Load application specs from datasource
func (m APIManager) getAPISpecs() []*APISpec {
	log.Debug("Using App Configuration from Mongo DB")
	return APILoader.LoadDefinitionsFromDatastore(m.accessor.Session)
}

//validateProxy validates proxy data
func (m APIManager) validateProxy(proxy Proxy) bool {
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
