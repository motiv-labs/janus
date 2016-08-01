package main

import (
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/etcinit/speedbump"
	"github.com/hellofresh/api-gateway/storage"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/redis.v3"
	"os"
	"github.com/gin-gonic/gin"
)

var APILoader = APIDefinitionLoader{}
var config = Specification{}

//loadConfigEnv loads environment variables
func loadConfigEnv() Specification {
	err := envconfig.Process("", &config)

	if err != nil {
		log.Fatal(err.Error())
	}

	return config
}

// initializeDatabase initializes a DB connection
func initializeDatabase() *storage.DatabaseAccessor {
	accessor, err := storage.NewServer(config.Database)
	if err != nil {
		log.Fatal(err)
	}

	return accessor
}

// initializeRedis initializes a Redis connection
func initializeRedis() *redis.Client {
	log.Infof("Trying to connect to redis instance: %s", config.Storage.DSN)
	return redis.NewClient(&redis.Options{
		Addr:     config.Storage.DSN,
		Password: config.Storage.Password,
		DB:       config.Storage.Database,
	})
}

//loadAPIEndpoints register api endpoints
func loadAPIEndpoints(router *gin.Engine) {
	log.Debug("Loading API Endpoints")

	handler := AppsAPI{}

	group := router.Group("/apis")
	{
		group.GET("/", handler.Get())
		group.POST("/", handler.Post())
		group.GET("/:id", handler.GetBy())
		group.PUT("/:id", handler.PutBy())
		group.DELETE("/:id", handler.DeleteBy())
	}
}

//getAPISpecs Load application specs from source
func getAPISpecs(accessor *storage.DatabaseAccessor) []*APISpec {
	log.Debug("Using App Configuration from Mongo DB")
	return APILoader.LoadDefinitionsFromDatastore(accessor.Session)
}

//loadApps load application middleware
func loadApps(router *gin.Engine, apiSpecs []*APISpec, redisClient *redis.Client, accessor *storage.DatabaseAccessor) {
	log.Debug("Loading API configurations")

	proxyRegister := &ProxyRegister{router}

	for _, referenceSpec := range apiSpecs {
		var skip bool

		//Define fields to log
		logger := createContextualLogger(referenceSpec)

		//Validates the proxy
		skip = validateProxy(referenceSpec.Proxy)

		if false == referenceSpec.Active {
			logger.Info("API is not active, skiping...")
			skip = false
		}

		if skip {
			cb := NewCircuitBreaker(referenceSpec)

			hasher := speedbump.PerSecondHasher{}
			limit := referenceSpec.RateLimit.Limit
			limiter := speedbump.NewLimiter(redisClient, hasher, limit)

			mw := &Middleware{referenceSpec, logger}

			router.Use(CreateMiddleware(&Database{mw, accessor}))

			var beforeHandlers = []gin.HandlerFunc{
				CreateMiddleware(&RateLimitMiddleware{mw, limiter, hasher, limit}),
			}

			if referenceSpec.UseOauth2 {
				logger.Debug("Loading OAuth Manager")
				referenceSpec.OAuthManager = &OAuthManager{redisClient}
				addOAuthHandlers(proxyRegister, mw, cb, logger)
				beforeHandlers = append(beforeHandlers, CreateMiddleware(Oauth2KeyExists{mw}))
				logger.Debug("Done loading OAuth Manager")
			}

			proxyRegister.Register(referenceSpec.Proxy, cb, beforeHandlers, nil)
			logger.Debug("Proxy registered")
		} else {
			logger.Error("Listen path is empty, skipping...")
		}
	}
}

//addOAuthHandlers loads configured oauth endpoints
func addOAuthHandlers(proxyRegister *ProxyRegister, mw *Middleware, cb *ExtendedCircuitBreakerMeta, logger *Logger) {
	logger.Info("Loading oauth configuration")
	var proxies []Proxy
	var handlers []gin.HandlerFunc

	oauthMeta := mw.Spec.Oauth2Meta

	//oauth proxy
	logger.Debug("Registering authorize endpoint")
	authorizeProxy := oauthMeta.OauthEndpoints.Authorize
	if validateProxy(authorizeProxy) {
		proxies = append(proxies, authorizeProxy)
	} else {
		logger.Debug("No authorize endpoint")
	}

	logger.Debug("Registering token endpoint")
	tokenProxy := oauthMeta.OauthEndpoints.Token
	if validateProxy(tokenProxy) {
		proxies = append(proxies, tokenProxy)
	} else {
		logger.Debug("No token endpoint")
	}

	logger.Debug("Registering info endpoint")
	infoProxy := oauthMeta.OauthEndpoints.Info
	if validateProxy(infoProxy) {
		proxies = append(proxies, infoProxy)
	} else {
		logger.Debug("No info endpoint")
	}

	logger.Debug("Registering create client endpoint")
	createProxy := oauthMeta.OauthClientEndpoints.Create
	if validateProxy(createProxy) {
		proxies = append(proxies, createProxy)
	} else {
		logger.Debug("No client create endpoint")
	}

	logger.Debug("Registering remove client endpoint")
	removeProxy := oauthMeta.OauthClientEndpoints.Remove
	if validateProxy(removeProxy) {
		proxies = append(proxies, removeProxy)
	} else {
		logger.Debug("No client remove endpoint")
	}

	handlers = append(handlers, CreateMiddleware(OAuthMiddleware{mw}))
	proxyRegister.registerMany(proxies, cb, nil, handlers)
}

//validateProxy validates proxy data
func validateProxy(proxy Proxy) bool {
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

func main() {
	loadConfigEnv()
	log.SetOutput(os.Stderr)

	if config.Debug {
		log.SetLevel(log.DebugLevel)
		gin.SetMode(gin.DebugMode)
	} else {
		log.SetLevel(log.InfoLevel)
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	accessor := initializeDatabase()
	defer accessor.Close()

	redisStorage := initializeRedis()
	defer redisStorage.Close()

	specs := getAPISpecs(accessor)
	loadApps(router, specs, redisStorage, accessor)
	loadAPIEndpoints(router)

	router.Run(fmt.Sprintf(":%v", config.Port))
}
