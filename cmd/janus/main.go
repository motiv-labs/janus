package main

import (
	"crypto/tls"
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/hellofresh/janus/pkg/api"
	"github.com/hellofresh/janus/pkg/jwt"
	"github.com/hellofresh/janus/pkg/loader"
	"github.com/hellofresh/janus/pkg/middleware"
	"github.com/hellofresh/janus/pkg/oauth"
	"github.com/hellofresh/janus/pkg/proxy"
	"github.com/hellofresh/janus/pkg/router"
)

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
	defer accessor.Close()
	defer statsdClient.Close()

	http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = globalConfig.MaxIdleConnsPerHost
	if globalConfig.InsecureSkipVerify {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	router := router.NewHttpTreeMuxRouter()
	router.Use(middleware.NewLogger(globalConfig.Debug).Handler, middleware.NewRecovery(RecoveryHandler).Handler, middleware.NewMongoDB(accessor).Handler)

	manager := &oauth.Manager{Storage: storage}
	transport := oauth.NewAwareTransport(http.DefaultTransport, manager)
	registerChan := proxy.NewRegisterChan(router, transport)

	changeTracker := loader.NewTracker()
	apiLoader := api.NewLoader(registerChan, storage, accessor, manager, globalConfig.Debug)
	apiLoader.Load()
	apiLoader.ListenToChanges(changeTracker)

	oauthLoader := oauth.NewLoader(registerChan, accessor, globalConfig.Debug)
	oauthLoader.Load()
	oauthLoader.ListenToChanges(changeTracker)

	authConfig := jwt.NewConfig(globalConfig.Credentials)
	authMiddleware := jwt.NewMiddleware(authConfig)

	// Home endpoint for the gateway
	router.GET("/", Home(globalConfig.Application))
	loadAuthEndpoints(router, authMiddleware)
	loadAPIEndpoints(router, authMiddleware, changeTracker)
	loadOAuthEndpoints(router, authMiddleware, changeTracker)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", globalConfig.Port), router))
}
