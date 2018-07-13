package web

import (
	"fmt"
	"net/http"
	"net/http/pprof"

	chiMiddleware "github.com/go-chi/chi/middleware"
	"github.com/hellofresh/janus/pkg/api"
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
	Port              int
	Credentials       config.Credentials
	TLS               config.TLS
	ConfigurationChan chan api.ConfigurationMessage
	apiHandler        *APIHandler
	profilingEnabled  bool
	profilingPublic   bool
}

// New creates a new web server
func New(opts ...Option) *Server {
	cfgChan := make(chan api.ConfigurationMessage)
	s := Server{
		ConfigurationChan: cfgChan,
		apiHandler:        NewAPIHandler(cfgChan),
	}

	for _, opt := range opts {
		opt(&s)
	}

	return &s
}

// Start creates a router and serves requests async
func (s *Server) Start() error {
	log.Info("Janus Admin API starting...")
	router.DefaultOptions.NotFoundHandler = httpErrors.NotFound
	r := router.NewChiRouterWithOptions(router.DefaultOptions)
	go s.listenAndServe(r)

	s.AddRoutes(r)
	plugin.EmitEvent(plugin.AdminAPIStartupEvent, plugin.OnAdminAPIStartup{Router: r})

	return nil
}

// Stop stops the server
func (s *Server) Stop() {
	close(s.ConfigurationChan)
}

// AddRoutes adds the admin routes
func (s *Server) AddRoutes(r router.Router) {
	// create authentication for Janus
	guard := jwt.NewGuard(s.Credentials)
	r.Use(
		chiMiddleware.StripSlashes,
		chiMiddleware.DefaultCompress,
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

	s.addInternalPublicRoutes(r)
	s.addInternalAuthRoutes(r, guard)
	s.addInternalRoutes(r, guard)
}

func (s *Server) addInternalPublicRoutes(r router.Router) {
	r.GET("/", Home())
	r.GET("/status", NewOverviewHandler(s.apiHandler.Cfgs))
	r.GET("/status/{name}", NewStatusHandler(s.apiHandler.Cfgs))
}

func (s *Server) addInternalAuthRoutes(r router.Router, guard jwt.Guard) {
	handlers := jwt.Handler{Guard: guard}
	r.POST("/login", handlers.Login(s.Credentials))
	authGroup := r.Group("/auth")
	{
		authGroup.GET("/refresh_token", handlers.Refresh())
	}
}

func (s *Server) addInternalRoutes(r router.Router, guard jwt.Guard) {
	log.Debug("Loading API Endpoints")

	// APIs endpoints
	groupAPI := r.Group("/apis")
	groupAPI.Use(jwt.NewMiddleware(guard).Handler)
	{
		groupAPI.GET("/", s.apiHandler.Get())
		groupAPI.GET("/{name}", s.apiHandler.GetBy())
		groupAPI.POST("/", s.apiHandler.Post())
		groupAPI.PUT("/{name}", s.apiHandler.PutBy())
		groupAPI.DELETE("/{name}", s.apiHandler.DeleteBy())
	}

	if s.profilingEnabled {
		groupProfiler := r.Group("/debug/pprof")
		if !s.profilingPublic {
			groupProfiler.Use(jwt.NewMiddleware(guard).Handler)
		}
		{
			groupProfiler.GET("/*", pprof.Index)
			groupProfiler.GET("/cmdline", pprof.Cmdline)
			groupProfiler.GET("/profile", pprof.Profile)
			groupProfiler.GET("/symbol", pprof.Symbol)
			groupProfiler.GET("/trace", pprof.Trace)
		}
	}
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
