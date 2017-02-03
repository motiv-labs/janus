package main

import (
	"fmt"
	"net/http"

	"github.com/NYTimes/gziphandler"
	log "github.com/Sirupsen/logrus"
	"github.com/hellofresh/janus/pkg/api"
	"github.com/hellofresh/janus/pkg/jwt"
	"github.com/hellofresh/janus/pkg/middleware"
	"github.com/hellofresh/janus/pkg/oauth"
	"github.com/hellofresh/janus/pkg/proxy"
	"github.com/hellofresh/janus/pkg/router"
	"github.com/hellofresh/janus/pkg/stats"
)

//loadAPIEndpoints register api endpoints
func loadAPIEndpoints(router router.Router, authMiddleware *jwt.Middleware) {
	log.Debug("Loading API Endpoints")

	// Apis endpoints
	handler := api.NewController()
	group := router.Group("/apis")
	group.Use(authMiddleware.Handler, gziphandler.GzipHandler)
	{
		group.GET("", handler.Get())
		group.POST("", handler.Post())
		group.GET("/:id", handler.GetBy())
		group.PUT("/:id", handler.PutBy())
		group.DELETE("/:id", handler.DeleteBy())
	}
}

//loadOAuthEndpoints register api endpoints
func loadOAuthEndpoints(router router.Router, authMiddleware *jwt.Middleware) {
	log.Debug("Loading OAuth Endpoints")

	// Oauth servers endpoints
	oAuthHandler := oauth.NewController()
	oauthGroup := router.Group("/oauth/servers")
	oauthGroup.Use(authMiddleware.Handler, gziphandler.GzipHandler)
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
	defer accessor.Close()
	defer statsdClient.Close()

	statsClient := stats.NewStatsClient(statsdClient)

	// create router
	r := router.NewHttpTreeMuxRouter()
	r.Use(
		middleware.NewStats(statsClient).Handler,
		middleware.NewLogger(globalConfig.Debug).Handler,
		middleware.NewRecovery(RecoveryHandler).Handler,
		middleware.NewMongoDB(accessor).Handler,
	)

	// create the proxy
	oAuthServersRepo, err := oauth.NewMongoRepository(accessor.Session.DB(""))
	if err != nil {
		log.Panic(err)
	}
	manager := &oauth.Manager{Storage: storage}
	transport := oauth.NewAwareTransport(manager, oAuthServersRepo, statsClient)
	p := proxy.WithParams(proxy.Params{
		Transport:              transport,
		FlushInterval:          globalConfig.BackendFlushInterval,
		IdleConnectionsPerHost: globalConfig.MaxIdleConnsPerHost,
		CloseIdleConnsPeriod:   globalConfig.CloseIdleConnsPeriod,
		InsecureSkipVerify:     globalConfig.InsecureSkipVerify,
	})
	defer p.Close()

	// create proxy register
	register := proxy.NewRegister(r, p)
	apiLoader := api.NewLoader(register, storage, accessor, manager, globalConfig.Debug)
	apiLoader.Load()

	oauthLoader := oauth.NewLoader(register, accessor, globalConfig.Debug)
	oauthLoader.Load()

	// create authentication for Janus
	authConfig := jwt.NewConfig(globalConfig.Credentials)
	authMiddleware := jwt.NewMiddleware(authConfig)

	// create endpoints
	r.GET("/", Home(globalConfig.Application))
	r.GET("/status", Heartbeat())

	loadAuthEndpoints(r, authMiddleware)
	loadAPIEndpoints(r, authMiddleware)
	loadOAuthEndpoints(r, authMiddleware)

	log.Fatal(listenAndServe(r))
}

func listenAndServe(handler http.Handler) error {
	address := fmt.Sprintf(":%v", globalConfig.Port)
	log.Infof("Listening on %v", address)
	if globalConfig.IsHTTPS() {
		return http.ListenAndServeTLS(address, globalConfig.CertPathTLS, globalConfig.KeyPathTLS, handler)
	}

	log.Infof("certPathTLS or keyPathTLS not found, defaulting to HTTP")
	return http.ListenAndServe(address, handler)
}
