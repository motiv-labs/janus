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

func loadAPIEndpoints(proxyRegister *ProxyRegister) {
	log.Debug("Loading API Endpoints")
	group := iris.Party("/apis")
	{
		appHandler := AppsAPI{proxyRegister: proxyRegister}
		group.API("/", appHandler)
	}
}

func getAPISpecs(accessor *storage.DatabaseAccessor) []*APISpec {
	log.Debug("Using App Configuration from Mongo DB")
	return APILoader.LoadDefinitionsFromDatastore(accessor.Session)
}

func loadApps(proxyRegister *ProxyRegister, apiSpecs []*APISpec, redisClient *redis.Client, accessor *storage.DatabaseAccessor) {
	log.Debug("Loading API configurations")

	for _, referenceSpec := range apiSpecs {
		skip := validateProxy(referenceSpec.Proxy)
		cb := NewCircuitBreaker(referenceSpec)

		if referenceSpec.UseOauth2 {
			loadOAuth(proxyRegister, referenceSpec.Oauth2Meta, cb)
		}

		if skip {
			proxyRegister.Register(referenceSpec.Proxy, cb)

			hasher := speedbump.PerSecondHasher{}
			limit := referenceSpec.RateLimit.Limit
			limiter := speedbump.NewLimiter(redisClient, hasher, limit)

			mw := &Middleware{referenceSpec}
			CreateMiddleware(&Database{mw, accessor}, mw)
			CreateMiddleware(&RateLimitMiddleware{mw, limiter, hasher, limit}, mw)

			log.Debug("Proxy registered")
		} else {
			log.Error("Listen path is empty, skipping API ID: ", referenceSpec.ID)
		}
	}
}

func loadOAuth(proxyRegister *ProxyRegister, oauthMeta Oauth2Meta, cb *ExtendedCircuitBreakerMeta) {
	log.Debug("Loading oauth configuration")
	var proxies []Proxy

	//oauth proxy
	log.Debug("Registering authorize endpoint")
	proxies = append(proxies, oauthMeta.OauthEndpoints.Authorize)

	log.Debug("Registering token endpoint")
	proxies = append(proxies, oauthMeta.OauthEndpoints.Token)

	log.Debug("Registering info endpoint")
	proxies = append(proxies, oauthMeta.OauthEndpoints.Info)

	log.Debug("Registering create client endpoint")
	proxies = append(proxies, oauthMeta.OauthClientEndpoints.Create)

	log.Debug("Registering remove client endpoint")
	proxies = append(proxies, oauthMeta.OauthClientEndpoints.Remove)

	proxyRegister.registerMany(proxies, cb)
}

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

	proxyRegister := NewProxyRegister()

	specs := getAPISpecs(accessor)
	loadApps(proxyRegister, specs, redisStorage, accessor)
	loadAPIEndpoints(proxyRegister)

	iris.Listen(fmt.Sprintf(":%v", config.Port))
}
