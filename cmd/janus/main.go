package main

import (
	"fmt"
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/hellofresh/janus/api"
	"github.com/hellofresh/janus/config"
	"github.com/hellofresh/janus/jwt"
	"github.com/hellofresh/janus/middleware"
	"github.com/hellofresh/janus/oauth"
	"github.com/hellofresh/janus/proxy"
	"github.com/hellofresh/janus/router"
	statsd "gopkg.in/alexcesaro/statsd.v2"
	"gopkg.in/redis.v3"
)

// initLogger initializes the logger config
func initLogger(config *config.Specification) {
	log.SetOutput(os.Stderr)
	// log.SetFormatter(&logstash.LogstashFormatter{Type: config.Application.Name})

	if config.Debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
}

// initializeDatabase initializes a DB connection
func initializeDatabase(dsn string) *middleware.DatabaseAccessor {
	accessor, err := middleware.InitDB(dsn)
	if err != nil {
		log.Fatalf("Couldn't connect to the mongodb database: %s", err.Error())
	}

	return accessor
}

// initializeRedis initializes a Redis connection
func initializeRedis(dsn string) *redis.Client {
	log.Debugf("Trying to connect to redis instance: %s", dsn)
	return redis.NewClient(&redis.Options{
		Addr: dsn,
	})
}

// Initializes new StatsD client if it enabled
func initializeStatsd(dsn, prefix string) *statsd.Client {
	var options []statsd.Option

	log.Debugf("Trying to connect to statsd instance: %s", dsn)
	if len(dsn) == 0 {
		log.Debug("Statsd DSN not provided, client will be muted")
		options = append(options, statsd.Mute(true))
	} else {
		options = append(options, statsd.Address(dsn))
	}

	if len(prefix) > 0 {
		options = append(options, statsd.Prefix(prefix))
	}

	client, err := statsd.New(options...)

	if err != nil {
		log.WithError(err).
			WithFields(log.Fields{
				"dsn":    dsn,
				"prefix": prefix,
			}).Warning("An error occurred while connecting to StatsD. Client will be muted.")
	}

	return client
}

//loadAPIEndpoints register api endpoints
func loadAPIEndpoints(router router.Router, loader *api.Loader, authMiddleware *jwt.Middleware, config *config.Specification) {
	log.Debug("Loading API Endpoints")

	// Home endpoint for the gateway
	router.GET("/", Home(config.Application))

	// Apis endpoints
	handler := api.API{loader}
	group := router.Group("/apis")
	group.Use(authMiddleware.Handler)
	{
		group.GET("", handler.Get())
		group.POST("", handler.Post())
		group.GET("/:id", handler.GetBy())
		group.PUT("/:id", handler.PutBy())
		group.DELETE("/:id", handler.DeleteBy())
	}
}

//loadOAuthEndpoints register api endpoints
func loadOAuthEndpoints(router router.Router, loader *oauth.Loader, authMiddleware *jwt.Middleware) {
	log.Debug("Loading OAuth Endpoints")

	// Oauth servers endpoints
	oAuthHandler := oauth.API{loader}
	oauthGroup := router.Group("/oauth/servers")
	oauthGroup.Use(authMiddleware.Handler)
	{
		oauthGroup.GET("", oAuthHandler.Get())
		oauthGroup.POST("", oAuthHandler.Post())
		oauthGroup.GET("/:id", oAuthHandler.GetBy())
		oauthGroup.PUT("/:id", oAuthHandler.PutBy())
		oauthGroup.DELETE("/:id", oAuthHandler.DeleteBy())
	}
}

func loadAuthEndpoints(router router.Router, authMiddleware *jwt.Middleware) {
	log.Debug("Loading Auth Endpoints")

	handlers := jwt.Handler{Config: authMiddleware.Config}
	router.POST("/login", handlers.Login())
	authGroup := router.Group("/auth")
	authGroup.Use(authMiddleware.Handler)
	{
		authGroup.GET("/refresh_token", handlers.Refresh())
	}
}

func main() {
	config, err := config.LoadEnv()
	if nil != err {
		log.Panic(err.Error())
	}
	initLogger(config)

	router := router.NewHttpTreeMuxRouter()
	accessor := initializeDatabase(config.DatabaseDSN)
	router.Use(middleware.NewLogger(config.Debug).Handler, middleware.NewRecovery(RecoveryHandler).Handler, middleware.NewMongoDB(accessor).Handler)

	redisStorage := initializeRedis(config.StorageDSN)
	defer redisStorage.Close()

	statsdClient := initializeStatsd(config.StatsdDSN, config.StatsdPrefix)
	defer statsdClient.Close()

	manager := &oauth.Manager{redisStorage}
	transport := oauth.NewAwareTransport(http.DefaultTransport, manager)
	proxyRegister := &proxy.Register{Router: router, Transport: transport}

	apiLoader := api.NewLoader(router, redisStorage, accessor, proxyRegister, manager, config.Debug)
	apiLoader.Load()

	oauthLoader := oauth.NewLoader(router, accessor, proxyRegister, config.Debug)
	oauthLoader.Load()

	authConfig := jwt.NewConfig(config.Credentials)
	authMiddleware := jwt.NewMiddleware(authConfig)

	loadAuthEndpoints(router, authMiddleware)
	loadAPIEndpoints(router, apiLoader, authMiddleware, config)
	loadOAuthEndpoints(router, oauthLoader, authMiddleware)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", config.Port), router))
}
