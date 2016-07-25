package main

import (
	"github.com/kelseyhightower/envconfig"
	log "github.com/Sirupsen/logrus"
	"fmt"
	"github.com/kataras/iris"
	"github.com/hellofresh/api-gateway/storage"
	"strings"
	"gopkg.in/redis.v3"
	"github.com/etcinit/speedbump"
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
func loadAPIEndpoints() {
	log.Debug("Loading API Endpoints")
	group := iris.Party("/apis")
	{
		appHandler := AppsAPI{}
		group.API("/", appHandler)
	}
}

//getAPISpecs Load application specs from source
func getAPISpecs(accessor *storage.DatabaseAccessor) []*APISpec {
	log.Debug("Using App Configuration from Mongo DB")
	return APILoader.LoadDefinitionsFromDatastore(accessor.Session)
}

//loadApps load application middleware
func loadApps(apiSpecs []*APISpec, redisClient *redis.Client, accessor *storage.DatabaseAccessor) {
	log.Debug("Loading API configurations")

	proxyRegister := NewProxyRegister()

	for _, referenceSpec := range apiSpecs {
		skip := validateProxy(referenceSpec.Proxy)
		cb := NewCircuitBreaker(referenceSpec)

		if skip {
			proxyRegister.Register(referenceSpec.Proxy, cb)

			hasher := speedbump.PerSecondHasher{}
			limit := referenceSpec.RateLimit.Limit
			limiter := speedbump.NewLimiter(redisClient, hasher, limit)

			mw := &Middleware{referenceSpec}
			CreateMiddleware(&Database{mw, accessor}, mw)
			CreateMiddleware(&RateLimitMiddleware{mw, limiter, hasher, limit}, mw)

			if referenceSpec.UseOauth2 {
				log.Debug("Loading OAuth Manager")
				referenceSpec.OAuthManager = &OAuthManager{redisClient}
				addOAuthHandlers(proxyRegister, referenceSpec, cb)
				CreateMiddleware(Oauth2KeyExists{mw}, mw)
				log.Debug("Done loading OAuth Manager")
			}

			log.Debug("Proxy registered")
		} else {
			log.Error("Listen path is empty, skipping API ID: ", referenceSpec.ID)
		}
	}
}

//addOAuthHandlers loads configured oauth endpoints
func addOAuthHandlers(proxyRegister *ProxyRegister, spec *APISpec, cb *ExtendedCircuitBreakerMeta) {
	log.Info("Loading oauth configuration")
	var proxies []Proxy
	oauthMeta := spec.Oauth2Meta

	//oauth proxy
	log.Debug("Registering authorize endpoint")
	authorizeProxy := oauthMeta.OauthEndpoints.Authorize
	if validateProxy(authorizeProxy) {
		proxies = append(proxies, authorizeProxy)
	} else {
		log.Debug("No authorize endpoint")
	}

	log.Debug("Registering token endpoint")
	tokenProxy := oauthMeta.OauthEndpoints.Token
	if validateProxy(tokenProxy) {
		proxies = append(proxies, tokenProxy)
	} else {
		log.Debug("No token endpoint")
	}

	log.Debug("Registering info endpoint")
	infoProxy := oauthMeta.OauthEndpoints.Info
	if validateProxy(infoProxy) {
		proxies = append(proxies, infoProxy)
	} else {
		log.Debug("No info endpoint")
	}

	log.Debug("Registering create client endpoint")
	createProxy := oauthMeta.OauthClientEndpoints.Create
	if validateProxy(createProxy) {
		proxies = append(proxies, createProxy)
	} else {
		log.Debug("No client create endpoint")
	}

	log.Debug("Registering remove client endpoint")
	removeProxy := oauthMeta.OauthClientEndpoints.Remove
	if validateProxy(removeProxy) {
		proxies = append(proxies, removeProxy)
	} else {
		log.Debug("No client remove endpoint")
	}

	proxyRegister.registerMany(proxies, cb, OAuthHandler{spec})
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
	log.SetLevel(log.DebugLevel)
	loadConfigEnv()

	accessor := initializeDatabase()
	defer accessor.Close()

	redisStorage := initializeRedis()
	defer redisStorage.Close()

	specs := getAPISpecs(accessor)
	loadApps(specs, redisStorage, accessor)
	loadAPIEndpoints()

	iris.Listen(fmt.Sprintf(":%v", config.Port))
}
