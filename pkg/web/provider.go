package web

import (
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/hellofresh/janus/pkg/api"
	"github.com/hellofresh/janus/pkg/config"
	"github.com/hellofresh/janus/pkg/jwt"
	"github.com/hellofresh/janus/pkg/middleware"
	"github.com/hellofresh/janus/pkg/notifier"
	"github.com/hellofresh/janus/pkg/oauth"
	"github.com/hellofresh/janus/pkg/router"
	chimiddleware "github.com/pressly/chi/middleware"
	"github.com/rs/cors"
)

// Provider is a provider.Provider implementation that provides the REST API.
type Provider struct {
	Port     int    `description:"Web administration port"`
	CertFile string `description:"SSL certificate"`
	KeyFile  string `description:"SSL certificate"`
	ReadOnly bool   `description:"Enable read only API"`
	Cred     config.Credentials
	APIRepo  api.Repository
	AuthRepo oauth.Repository
	Notifier notifier.Notifier
}

// Provide executes the provider functionality
// This is normally the entry point of the
// provider.
func (p *Provider) Provide(version string) error {
	r := router.NewChiRouter()

	// create authentication for Janus
	authConfig := jwt.NewConfig(p.Cred)
	authMiddleware := jwt.NewMiddleware(authConfig)
	r.Use(
		chimiddleware.StripSlashes,
		chimiddleware.DefaultCompress,
		middleware.NewLogger().Handler,
		middleware.NewRecovery(RecoveryHandler).Handler,
		middleware.NewOpenTracing(p.IsHTTPS()).Handler,
		cors.New(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedHeaders:   []string{"*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
			AllowCredentials: true,
		}).Handler,
	)

	// create endpoints
	r.GET("/", Home(version))
	r.GET("/status", Heartbeat())

	handlers := jwt.Handler{Config: authConfig}
	r.POST("/login", handlers.Login())
	authGroup := r.Group("/auth")
	{
		authGroup.GET("/refresh_token", handlers.Refresh())
	}

	p.loadAPIEndpoints(r, authMiddleware.Handler)
	p.loadOAuthEndpoints(r)

	go func() {
		log.Fatal(p.listenAndServe(r))
	}()
	return nil
}

func (p *Provider) listenAndServe(handler http.Handler) error {
	address := fmt.Sprintf(":%v", p.Port)
	log.WithField("address", address).Info("Listening on")
	log.Info("Janus Admin API started")
	if p.IsHTTPS() {
		return http.ListenAndServeTLS(address, p.CertFile, p.KeyFile, handler)
	}

	log.Info("Certificate and certificate key were not found, defaulting to HTTP")
	return http.ListenAndServe(address, handler)
}

//loadAPIEndpoints register api endpoints
func (p *Provider) loadAPIEndpoints(router router.Router, handlers ...router.Constructor) {
	log.Debug("Loading API Endpoints")

	// Apis endpoints
	handler := api.NewController(p.APIRepo, p.Notifier)
	group := router.Group("/apis")
	{
		group.GET("/", handler.Get())
		group.GET("/:name", handler.GetBy())

		if false == p.ReadOnly {
			group.POST("/", handler.Post())
			group.PUT("/:name", handler.PutBy())
			group.DELETE("/:name", handler.DeleteBy())
		}
	}
}

//loadOAuthEndpoints register api endpoints
func (p *Provider) loadOAuthEndpoints(router router.Router, handlers ...router.Constructor) {
	log.Debug("Loading OAuth Endpoints")

	// Oauth servers endpoints
	oAuthHandler := oauth.NewController(p.AuthRepo, p.Notifier)
	oauthGroup := router.Group("/oauth/servers")
	{
		oauthGroup.GET("/", oAuthHandler.Get(), handlers...)
		oauthGroup.GET("/:name", oAuthHandler.GetBy(), handlers...)

		if false == p.ReadOnly {
			oauthGroup.POST("/", oAuthHandler.Post(), handlers...)
			oauthGroup.PUT("/:name", oAuthHandler.PutBy(), handlers...)
			oauthGroup.DELETE("/:name", oAuthHandler.DeleteBy(), handlers...)
		}
	}
}

// IsHTTPS checks if you have https enabled
func (p *Provider) IsHTTPS() bool {
	return len(p.CertFile) > 0 && len(p.KeyFile) > 0
}
