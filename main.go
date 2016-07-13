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

// Specification for basic configurations
type Specification struct {
	storage.Database
	Port  int                 `envconfig:"PORT"`
	Debug bool                `envconfig:"DEBUG"`
	RedisSpecification
}

type RedisSpecification struct {
	DSN      string                `envconfig:"REDIS_DSN"`
	Password string                `envconfig:"REDIS_PASSWORD"`
	DB       int64                 `envconfig:"REDIS_DB"`
}

func loadConfigEnv() Specification {
	var s Specification
	err := envconfig.Process("", &s)

	if err != nil {
		log.Fatal(err.Error())
	}

	return s
}

// initializeDatabase initializes a DB connection
func initializeDatabase(dbConfig storage.Database) *storage.DatabaseAccessor {
	accessor, err := storage.NewServer(dbConfig)
	if err != nil {
		log.Fatal(err)
	}

	//Use the middleware
	iris.Use(NewDatabase(*accessor))

	return accessor
}

// initializeRedis initializes a Redis connection
func initializeRedis(spec RedisSpecification) *redis.Client {
	log.Infof("Trying to connect to %s", spec.DSN)
	return redis.NewClient(&redis.Options{
		Addr:     spec.DSN,
		Password: spec.Password,
		DB:       spec.DB,
	})
}

func loadAPIEndpoints(proxyRegister *ProxyRegister) {
	group := iris.Party("/apis")
	{
		appHandler := AppsAPI{proxyRegister: proxyRegister}
		group.API("/", appHandler)
	}
}

func getAPISpecs(accessor *storage.DatabaseAccessor, dbConfig storage.Database) []*APIDefinition {
	var APISpecs []*APIDefinition

	log.Debug("Using App Configuration from Mongo DB")
	APISpecs = APILoader.LoadDefinitionsFromDatastore(accessor.Session, dbConfig)

	return APISpecs;
}

func loadApps(apiSpecs []*APIDefinition, redisClient *redis.Client) {
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
			var middleware []iris.Handler

			if referenceSpec.RateLimit.Enabled {
				rateLimit := NewRateLimitMiddleware(redisClient, speedbump.PerSecondHasher{}, referenceSpec.RateLimit.Limit)
				middleware = append(middleware, rateLimit)
			}

			// Circuit Breaker Middleware
			cb := NewCircuitBreaker(referenceSpec)
			proxyRegister := NewProxyRegister()
			proxyRegister.Register(referenceSpec.Proxy, cb, middleware...)
			log.Debug("Proxy registered")
		}
	}
}

func main() {
	log.SetLevel(log.DebugLevel)

	s := loadConfigEnv()
	accessor := initializeDatabase(s.Database)
	defer accessor.Close()

	redis := initializeRedis(s.RedisSpecification)
	defer redis.Close()

	specs := getAPISpecs(accessor, s.Database)
	loadApps(specs, redis)

	proxyRegister := NewProxyRegister()
	loadAPIEndpoints(proxyRegister)

	iris.Listen(fmt.Sprintf(":%v", s.Port))
}
