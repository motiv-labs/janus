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
	"github.com/hellofresh/janus/pkg/plugin"
	"github.com/hellofresh/janus/pkg/router"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
)

// Server represents the web server
type Server struct {
	repo        api.Repository
	Port        int
	ReadOnly    bool
	Credentials config.Credentials
	TLS         config.TLS
}

// New creates a new web server
func New(repo api.Repository, opts ...Option) *Server {
	s := Server{repo: repo}

	for _, opt := range opts {
		opt(&s)
	}

	return &s
}

// Serve creates a router and serves requests async
func (s *Server) Serve() error {
	log.Info("Janus Admin API starting...")
	router.DefaultOptions.NotFoundHandler = httpErrors.NotFound
	r := router.NewChiRouterWithOptions(router.DefaultOptions)

	// create authentication for Janus
	guard := jwt.NewGuard(s.Credentials)
	r.Use(
		chimiddleware.StripSlashes,
		chimiddleware.DefaultCompress,
		middleware.NewLogger().Handler,
		middleware.NewRecovery(httpErrors.RecoveryHandler),
		middleware.NewOpenTracing(s.TLS.IsHTTPS()).Handler,
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
	r.GET("/status", checker.NewOverviewHandler(s.repo))
	r.GET("/status/{name}", checker.NewStatusHandler(s.repo))

	handlers := jwt.Handler{Guard: guard}
	r.POST("/login", handlers.Login(s.Credentials))
	authGroup := r.Group("/auth")
	{
		authGroup.GET("/refresh_token", handlers.Refresh())
	}

	s.loadAPIEndpoints(r, guard)
	plugin.EmitEvent(plugin.AdminAPIStartupEvent, plugin.OnAdminAPIStartup{Router: r})

	go func() {
		s.listenAndServe(r)
	}()

	return nil
}

func (s *Server) listenAndServe(handler http.Handler) error {
	address := fmt.Sprintf(":%v", s.Port)

	log.Info("Janus Admin API started")
	if s.TLS.IsHTTPS() {
		addressTLS := fmt.Sprintf(":%v", s.TLS.Port)
		if s.TLS.Redirect {
			go func() {
				log.WithField("address", address).Info("Listening HTTP redirects to HTTPS")
				log.Fatal(http.ListenAndServe(address, RedirectHTTPS(s.TLS.Port)))
			}()
		}

		log.WithField("address", addressTLS).Info("Listening HTTPS")
		return http.ListenAndServeTLS(addressTLS, s.TLS.CertFile, s.TLS.KeyFile, handler)
	}

	log.WithField("address", address).Info("Certificate and certificate key were not found, defaulting to HTTP")
	return http.ListenAndServe(address, handler)
}

//loadAPIEndpoints register api endpoints
func (s *Server) loadAPIEndpoints(r router.Router, guard jwt.Guard) {
	log.Debug("Loading API Endpoints")

	// Apis endpoints
	handler := NewAPIHandler(s.repo)
	group := r.Group("/apis")
	group.Use(jwt.NewMiddleware(guard).Handler)
	{
		group.GET("/", handler.Get())
		group.GET("/{name}", handler.GetBy())
		group.POST("/", handler.Post())
		group.PUT("/{name}", handler.PutBy())
		group.DELETE("/{name}", handler.DeleteBy())
	}
}
