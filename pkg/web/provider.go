package web

import (
	"fmt"
	"net/http"

	chimiddleware "github.com/go-chi/chi/middleware"
	"github.com/hellofresh/janus/pkg/api"
	"github.com/hellofresh/janus/pkg/checker"
	"github.com/hellofresh/janus/pkg/config"
	"github.com/hellofresh/janus/pkg/errors"
	"github.com/hellofresh/janus/pkg/jwt"
	"github.com/hellofresh/janus/pkg/middleware"
	"github.com/hellofresh/janus/pkg/notifier"
	"github.com/hellofresh/janus/pkg/oauth"
	"github.com/hellofresh/janus/pkg/router"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
)

// Provider is a provider.Provider implementation that provides the REST API.
type Provider struct {
	Port     int  `description:"Web administration port"`
	ReadOnly bool `description:"Enable read only API"`
	Cred     config.Credentials
	Notifier notifier.Notifier
	TLS      config.TLS
	APIRepo  api.Repository
	AuthRepo oauth.Repository
}

// Provide executes the provider functionality
// This is normally the entry point of the
// provider.
func (p *Provider) Provide(version string) error {
	log.Info("Janus Admin API starting...")

	router.DefaultOptions.NotFoundHandler = errors.NotFound
	r := router.NewChiRouterWithOptions(router.DefaultOptions)

	// create authentication for Janus
	guard := jwt.NewGuard(p.Cred)
	authMiddleware := jwt.NewMiddleware(guard)
	r.Use(
		chimiddleware.StripSlashes,
		chimiddleware.DefaultCompress,
		middleware.NewLogger().Handler,
		middleware.NewRecovery(errors.RecoveryHandler),
		middleware.NewOpenTracing(p.TLS.IsHTTPS()).Handler,
		cors.New(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedHeaders:   []string{"*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
			AllowCredentials: true,
		}).Handler,
	)

	// create endpoints
	r.GET("/", Home(version))
	// health checks
	r.GET("/status", checker.NewOverviewHandler(p.APIRepo))
	r.GET("/status/{name}", checker.NewStatusHandler(p.APIRepo))

	handlers := jwt.Handler{Guard: guard}
	r.POST("/login", handlers.Login())
	authGroup := r.Group("/auth")
	{
		authGroup.GET("/refresh_token", handlers.Refresh())
	}

	p.loadAPIEndpoints(r, authMiddleware.Handler)
	p.loadOAuthEndpoints(r, authMiddleware.Handler)

	go func() {
		p.listenAndServe(r)
	}()

	return nil
}

func (p *Provider) listenAndServe(handler http.Handler) error {
	address := fmt.Sprintf(":%v", p.Port)

	log.Info("Janus Admin API started")
	if p.TLS.IsHTTPS() {
		addressTLS := fmt.Sprintf(":%v", p.TLS.Port)
		if p.TLS.Redirect {
			go func() {
				log.WithField("address", address).Info("Listening HTTP redirects to HTTPS")
				log.Fatal(http.ListenAndServe(address, RedirectHTTPS(p.TLS.Port)))
			}()
		}

		log.WithField("address", addressTLS).Info("Listening HTTPS")
		return http.ListenAndServeTLS(addressTLS, p.TLS.CertFile, p.TLS.KeyFile, handler)
	}

	log.WithField("address", address).Info("Certificate and certificate key were not found, defaulting to HTTP")
	return http.ListenAndServe(address, handler)
}

//loadAPIEndpoints register api endpoints
func (p *Provider) loadAPIEndpoints(router router.Router, handlers ...router.Constructor) {
	log.Debug("Loading API Endpoints")

	// Apis endpoints
	handler := api.NewController(p.APIRepo, p.Notifier)
	group := router.Group("/apis")
	group.Use(handlers...)
	{
		group.GET("/", handler.Get())
		group.GET("/{name}", handler.GetBy())

		if false == p.ReadOnly {
			group.POST("/", handler.Post())
			group.PUT("/{name}", handler.PutBy())
			group.DELETE("/{name}", handler.DeleteBy())
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
		oauthGroup.GET("/{name}", oAuthHandler.GetBy(), handlers...)

		if false == p.ReadOnly {
			oauthGroup.POST("/", oAuthHandler.Post(), handlers...)
			oauthGroup.PUT("/{name}", oAuthHandler.PutBy(), handlers...)
			oauthGroup.DELETE("/{name}", oAuthHandler.DeleteBy(), handlers...)
		}
	}
}
