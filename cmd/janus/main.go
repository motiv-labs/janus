package main

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/fvbock/endless"
	"github.com/hellofresh/janus/api"
	"github.com/hellofresh/janus/config"
	"github.com/hellofresh/janus/jwt"
	"github.com/hellofresh/janus/loader"
	"github.com/hellofresh/janus/log"
	"github.com/hellofresh/janus/middleware"
	"github.com/hellofresh/janus/oauth"
	"github.com/hellofresh/janus/proxy"
	"github.com/hellofresh/janus/router"
	statsd "gopkg.in/alexcesaro/statsd.v2"
	"gopkg.in/redis.v3"
)

// initializeDatabase initializes a DB connection
func initializeDatabase(config config.Database) *middleware.DatabaseAccessor {
	accessor, err := middleware.InitDB(config.DSN)
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
func initializeStatsd(config config.Statsd) *statsd.Client {
	var options []statsd.Option

	log.Debugf("Trying to connect to statsd instance: %s", config.DSN)
	if len(config.DSN) == 0 {
		log.Debug("Statsd DSN not provided, client will be muted")
		options = append(options, statsd.Mute(true))
	} else {
		options = append(options, statsd.Address(config.DSN))
	}

	if len(config.Prefix) > 0 {
		options = append(options, statsd.Prefix(config.Prefix))
	}

	client, err := statsd.New(options...)

	if err != nil {
		log.WithError(err).
			WithFields(logrus.Fields{
				"dsn":    config.DSN,
				"prefix": config.Prefix,
			}).Warning("An error occurred while connecting to StatsD. Client will be muted.")
	}

	return client
}

//loadAPIEndpoints register api endpoints
func loadAPIEndpoints(router router.Router, authMiddleware *jwt.Middleware, changeTracker *loader.Tracker) {
	log.Debug("Loading API Endpoints")

	// Apis endpoints
	handler := api.NewController(changeTracker)
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
func loadOAuthEndpoints(router router.Router, authMiddleware *jwt.Middleware, changeTracker *loader.Tracker) {
	log.Debug("Loading OAuth Endpoints")

	// Oauth servers endpoints
	oAuthHandler := oauth.NewController(changeTracker)
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
	// load global configuration
	config, err := config.LoadEnv()
	if nil != err {
		log.Panic(err.Error())
	}

	http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = config.MaxIdleConnsPerHost
	if config.InsecureSkipVerify {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	// logging
	log.InitLog(config)

	accessor := initializeDatabase(config.Database)
	defer accessor.Close()

	redisStorage := initializeRedis(config.StorageDSN)
	defer redisStorage.Close()

	statsdClient := initializeStatsd(config.Statsd)
	defer statsdClient.Close()

	router := router.NewHttpTreeMuxRouter()
	router.Use(middleware.NewLogger(config.Debug).Handler, middleware.NewRecovery(RecoveryHandler).Handler, middleware.NewMongoDB(accessor).Handler)

	manager := &oauth.Manager{redisStorage}
	transport := oauth.NewAwareTransport(http.DefaultTransport, manager)
	registerChan := proxy.NewRegisterChan(router, transport)

	changeTracker := loader.NewTracker()
	apiLoader := api.NewLoader(registerChan, redisStorage, accessor, manager, config.Debug)
	apiLoader.Load()
	apiLoader.ListenToChanges(changeTracker)

	oauthLoader := oauth.NewLoader(registerChan, accessor, config.Debug)
	oauthLoader.Load()
	oauthLoader.ListenToChanges(changeTracker)

	authConfig := jwt.NewConfig(config.Credentials)
	authMiddleware := jwt.NewMiddleware(authConfig)

	// Home endpoint for the gateway
	router.GET("/", Home(config.Application))
	loadAuthEndpoints(router, authMiddleware)
	loadAPIEndpoints(router, authMiddleware, changeTracker)
	loadOAuthEndpoints(router, authMiddleware, changeTracker)

	log.Fatal(endless.ListenAndServe(fmt.Sprintf(":%v", config.Port), router))
}
