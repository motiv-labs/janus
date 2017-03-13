package web

import (
	"fmt"
	"net/http"

	"github.com/NYTimes/gziphandler"
	log "github.com/Sirupsen/logrus"
	"github.com/hellofresh/janus/pkg/api"
	"github.com/hellofresh/janus/pkg/config"
	"github.com/hellofresh/janus/pkg/jwt"
	"github.com/hellofresh/janus/pkg/middleware"
	"github.com/hellofresh/janus/pkg/oauth"
	"github.com/hellofresh/janus/pkg/router"
)

// Provider is a provider.Provider implementation that provides the REST API.
type Provider struct {
	Cred     config.Credentials
	APIRepo  api.APISpecRepository
	AuthRepo oauth.Repository
	Port     string
	CertFile string
	KeyFile  string
}

func (p *Provider) Provide() error {
	r := router.NewHttpTreeMuxRouter()

	// create authentication for Janus
	authConfig := jwt.NewConfig(p.Cred)
	authMiddleware := jwt.NewMiddleware(authConfig)
	r.Use(
		middleware.NewLogger().Handler,
		middleware.NewRecovery(RecoveryHandler).Handler,
		gziphandler.GzipHandler,
	)

	// create endpoints
	r.GET("/", Home())
	r.GET("/status", Heartbeat())
	handlers := jwt.Handler{Config: authConfig}
	r.POST("/login", handlers.Login())
	authGroup := r.Group("/auth")
	{
		authGroup.GET("/refresh_token", handlers.Refresh())
	}

	r.Use(authMiddleware.Handler)
	p.loadAPIEndpoints(r)
	p.loadOAuthEndpoints(r)

	go func() {
		log.Fatal(p.listenAndServe(r))
	}()
	return nil
}

func (p *Provider) listenAndServe(handler http.Handler) error {
	address := fmt.Sprintf(":%v", p.Port)
	log.Infof("Listening on %v", address)
	if len(p.CertFile) > 0 && len(p.KeyFile) > 0 {
		return http.ListenAndServeTLS(address, p.CertFile, p.KeyFile, handler)
	}

	log.Infof("certPathTLS or keyPathTLS not found, defaulting to HTTP")
	return http.ListenAndServe(address, handler)
}

//loadAPIEndpoints register api endpoints
func (p *Provider) loadAPIEndpoints(router router.Router) {
	log.Debug("Loading API Endpoints")

	// Apis endpoints
	handler := api.NewController(p.APIRepo)
	group := router.Group("/apis")
	{
		group.GET("", handler.Get())
		group.POST("", handler.Post())
		group.GET("/:slug", handler.GetBy())
		group.PUT("/:slug", handler.PutBy())
		group.DELETE("/:slug", handler.DeleteBy())
	}
}

//loadOAuthEndpoints register api endpoints
func (p *Provider) loadOAuthEndpoints(router router.Router) {
	log.Debug("Loading OAuth Endpoints")

	// Oauth servers endpoints
	oAuthHandler := oauth.NewController(p.AuthRepo)
	oauthGroup := router.Group("/oauth/servers")
	{
		oauthGroup.GET("", oAuthHandler.Get())
		oauthGroup.POST("", oAuthHandler.Post())
		oauthGroup.GET("/:slug", oAuthHandler.GetBy())
		oauthGroup.PUT("/:slug", oAuthHandler.PutBy())
		oauthGroup.DELETE("/:slug", oAuthHandler.DeleteBy())
	}
}
