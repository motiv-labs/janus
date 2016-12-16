package main

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/hellofresh/ginger-middleware/mongodb"
	"github.com/hellofresh/ginger-middleware/nice"
	"github.com/hellofresh/janus"
	"github.com/hellofresh/janus/config"
	statsd "gopkg.in/alexcesaro/statsd.v2"
	"gopkg.in/appleboy/gin-jwt.v2"
	"gopkg.in/redis.v3"
	"github.com/dimfeld/httptreemux"
	"github.com/urfave/negroni"
)

// initLogger initializes the logger config
func initLogger(config *config.Specification) {
	// log.SetFormatter(&logstash.LogstashFormatter{Type: config.Application.Name})

	if config.Debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
}

// initializeDatabase initializes a DB connection
func initializeDatabase(dsn string) *mongodb.DatabaseAccessor {
	accessor, err := mongodb.InitDB(dsn)
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
			}).
			Warning("An error occurred while connecting to StatsD. Client will be muted.")
	}

	return client
}

//loadDefaultEndpoints register api endpoints
func loadDefaultEndpoints(router *gin.Engine, apiManager *janus.APIManager, authMiddleware *jwt.GinJWTMiddleware, config *config.Specification) {
	log.Debug("Loading Default Endpoints")

	// Home endpoint for the gateway
	router.GET("/", janus.Home(config.Application))

	// Apis endpoints
	handler := janus.AppsAPI{apiManager}
	group := router.Group("/apis")
	group.Use(authMiddleware.MiddlewareFunc())
	{
		group.GET("", handler.Get())
		group.POST("", handler.Post())
		group.GET("/:id", handler.GetBy())
		group.PUT("/:id", handler.PutBy())
		group.DELETE("/:id", handler.DeleteBy())
	}

	// Oauth servers endpoints
	oAuthHandler := janus.OAuthAPI{}
	oauthGroup := router.Group("/oauth/servers")
	oauthGroup.Use(authMiddleware.MiddlewareFunc())
	{
		oauthGroup.GET("", oAuthHandler.Get())
		oauthGroup.POST("", oAuthHandler.Post())
		oauthGroup.GET("/:id", oAuthHandler.GetBy())
		oauthGroup.PUT("/:id", oAuthHandler.PutBy())
		oauthGroup.DELETE("/:id", oAuthHandler.DeleteBy())
	}
}

func loadAuthEndpoints(router *gin.Engine, authMiddleware *jwt.GinJWTMiddleware) {
	log.Debug("Loading Auth Endpoints")

	router.POST("/login", authMiddleware.LoginHandler)
	authGroup := router.Group("/auth")
	authGroup.Use(authMiddleware.MiddlewareFunc())
	{
		authGroup.GET("/refresh_token", authMiddleware.RefreshHandler)
	}
}

func main() {
	log.SetOutput(os.Stderr)

	config, err := config.LoadEnv()
	if nil != err {
		log.Panic(err.Error())
	}
	initLogger(config)

	router = janus.NewHttpTreeMuxRouter()
	n := negroni.New() // Includes some default middlewares
  	n.UseHandler(router)

	accessor := initializeDatabase(config.DatabaseDSN)
	n.Use(negroni.NewLogger(), middleware.Recovery(janus.RecoveryHandler), middleware.MongoSession(accessor))

	redisStorage := initializeRedis(config.StorageDSN)
	defer redisStorage.Close()

	statsdClient := initializeStatsd(config.StatsdDSN, config.StatsdPrefix)
	defer statsdClient.Close()

	apiManager := janus.NewAPIManager(router, redisStorage, accessor, statsdClient)
	apiManager.Load()

	authMiddleware := janus.NewJwt(&janus.Credentials{
		Secret:   config.Credentials.Secret,
		Username: config.Credentials.Username,
		Password: config.Credentials.Password,
	})

	loadAuthEndpoints(n, authMiddleware)
	loadDefaultEndpoints(n, apiManager, authMiddleware, config)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", config.Port), n)
}
