package web

import (
	"fmt"
	"net/http"

	chimiddleware "github.com/go-chi/chi/middleware"
	"github.com/hellofresh/janus/pkg/api"
	"github.com/hellofresh/janus/pkg/checker"
	"github.com/hellofresh/janus/pkg/config"
	httpErrors "github.com/hellofresh/janus/pkg/errors"
	"github.com/hellofresh/janus/pkg/jwt"
	"github.com/hellofresh/janus/pkg/middleware"
	"github.com/hellofresh/janus/pkg/notifier"
	"github.com/hellofresh/janus/pkg/plugin"
	"github.com/hellofresh/janus/pkg/router"
	"github.com/pkg/errors"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
)

func init() {
	plugin.RegisterEventHook(plugin.StartupEvent, onStartup)
}

func onStartup(event interface{}) error {
	e, ok := event.(plugin.OnStartup)
	if !ok {
		return errors.New("Could not convert event to startup type")
	}

	config := e.Config.Web

	repo, err := api.BuildRepository(e.Config.Database.DSN, e.MongoSession)
	if err != nil {
		return err
	}

	log.Info("Janus Admin API starting...")
	router.DefaultOptions.NotFoundHandler = httpErrors.NotFound
	r := router.NewChiRouterWithOptions(router.DefaultOptions)

	// create authentication for Janus
	guard := jwt.NewGuard(config.Credentials)
	r.Use(
		chimiddleware.StripSlashes,
		chimiddleware.DefaultCompress,
		middleware.NewLogger().Handler,
		middleware.NewRecovery(httpErrors.RecoveryHandler),
		middleware.NewOpenTracing(config.TLS.IsHTTPS()).Handler,
		cors.New(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedHeaders:   []string{"*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
			AllowCredentials: true,
		}).Handler,
	)

	// create endpoints
	r.GET("/", Home())
	// health checks
	r.GET("/status", checker.NewOverviewHandler(repo))
	r.GET("/status/{name}", checker.NewStatusHandler(repo))

	handlers := jwt.Handler{Guard: guard}
	r.POST("/login", handlers.Login(config.Credentials))
	authGroup := r.Group("/auth")
	{
		authGroup.GET("/refresh_token", handlers.Refresh())
	}

	loadAPIEndpoints(r, repo, e.Notifier, guard)
	plugin.EmitEvent(plugin.AdminAPIStartupEvent, plugin.OnAdminAPIStartup{Router: r})

	go func() {
		listenAndServe(config, r)
	}()

	return nil
}

func listenAndServe(config config.Web, handler http.Handler) error {
	address := fmt.Sprintf(":%v", config.Port)

	log.Info("Janus Admin API started")
	if config.TLS.IsHTTPS() {
		addressTLS := fmt.Sprintf(":%v", config.TLS.Port)
		if config.TLS.Redirect {
			go func() {
				log.WithField("address", address).Info("Listening HTTP redirects to HTTPS")
				log.Fatal(http.ListenAndServe(address, RedirectHTTPS(config.TLS.Port)))
			}()
		}

		log.WithField("address", addressTLS).Info("Listening HTTPS")
		return http.ListenAndServeTLS(addressTLS, config.TLS.CertFile, config.TLS.KeyFile, handler)
	}

	log.WithField("address", address).Info("Certificate and certificate key were not found, defaulting to HTTP")
	return http.ListenAndServe(address, handler)
}

//loadAPIEndpoints register api endpoints
func loadAPIEndpoints(router router.Router, repo api.Repository, ntf notifier.Notifier, guard jwt.Guard) {
	log.Debug("Loading API Endpoints")

	// Apis endpoints
	handler := api.NewController(repo, ntf)
	group := router.Group("/apis")
	group.Use(jwt.NewMiddleware(guard).Handler)
	{
		group.GET("/", handler.Get())
		group.GET("/{name}", handler.GetBy())
		group.POST("/", handler.Post())
		group.PUT("/{name}", handler.PutBy())
		group.DELETE("/{name}", handler.DeleteBy())
	}
}
