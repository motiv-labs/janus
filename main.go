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
	log.Infof("Trying to connect to %s", config.Storage.DSN)
	return redis.NewClient(&redis.Options{
		Addr:     config.Storage.DSN,
		Password: config.Storage.Password,
		DB:       config.Storage.Database,
	})
}

func loadAPIEndpoints(proxyRegister *ProxyRegister) {
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

func loadApps(apiSpecs []*APISpec, redisClient *redis.Client, accessor *storage.DatabaseAccessor) {
	log.Info("Loading API configurations.")

	for _, referenceSpec := range apiSpecs {
		var skip bool

		if referenceSpec.Proxy.ListenPath == "" {
			log.Error("Listen path is empty, skipping API ID: ", referenceSpec.ID)
			skip = true
		}

		if strings.Contains(referenceSpec.Proxy.ListenPath, " ") {
			log.Error("Listen path contains spaces, is invalid, skipping API ID: ", referenceSpec.ID)
			skip = true
		}

		if !skip {
			cb := NewCircuitBreaker(referenceSpec)
			proxyRegister := NewProxyRegister()
			proxyRegister.Register(referenceSpec.Proxy, cb)
			
			hasher := speedbump.PerSecondHasher{}
			limit := referenceSpec.RateLimit.Limit
			limiter := speedbump.NewLimiter(redisClient, hasher, limit)

			mw := &Middleware{referenceSpec}
			CreateMiddleware(&Database{mw, accessor}, mw)
			CreateMiddleware(&RateLimitMiddleware{mw, limiter, hasher, limit}, mw)
			CreateMiddleware(&OauthProxy{mw, proxyRegister}, mw)

			log.Debug("Proxy registered")
		}
	}
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

	proxyRegister := NewProxyRegister()
	loadAPIEndpoints(proxyRegister)

	iris.Listen(fmt.Sprintf(":%v", config.Port))
}
