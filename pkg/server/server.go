package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/cnf/structhash"
	"github.com/hellofresh/janus/pkg/api"
	"github.com/hellofresh/janus/pkg/config"
	"github.com/hellofresh/janus/pkg/errors"
	httpErrors "github.com/hellofresh/janus/pkg/errors"
	"github.com/hellofresh/janus/pkg/loader"
	"github.com/hellofresh/janus/pkg/middleware"
	"github.com/hellofresh/janus/pkg/plugin"
	"github.com/hellofresh/janus/pkg/proxy"
	"github.com/hellofresh/janus/pkg/router"
	"github.com/hellofresh/janus/pkg/web"
	"github.com/hellofresh/stats-go/client"
	log "github.com/sirupsen/logrus"
)

// Server is the Janus server
type Server struct {
	server       *http.Server
	repo         api.Repository
	register     *proxy.Register
	defLoader    *loader.APILoader
	started      bool
	version      string
	globalConfig *config.Specification
	statsClient  client.Client
}

// New creates a new instance of Server
func New(opts ...Option) *Server {
	s := Server{}

	for _, opt := range opts {
		opt(&s)
	}

	return &s
}

// StartWithContext starts the server and Stop/Close it when context is Done
func (s *Server) StartWithContext(ctx context.Context) error {
	go func() {
		defer s.Close()
		<-ctx.Done()
		log.Info("I have to go...")
		reqAcceptGraceTimeOut := time.Duration(s.globalConfig.GraceTimeOut)
		if reqAcceptGraceTimeOut > 0 {
			log.Infof("Waiting %s for incoming requests to cease", reqAcceptGraceTimeOut)
			time.Sleep(reqAcceptGraceTimeOut)
		}
		log.Info("Stopping server gracefully")
		s.Close()
	}()

	repo, err := api.BuildRepository(s.globalConfig.Database.DSN, s.globalConfig.Cluster.UpdateFrequency)
	s.repo = repo
	defer repo.Close()
	if err != nil {
		return errors.Wrap(err, "could not build a repository for the database")
	}

	ch := repo.Watch(ctx)
	go func() {
		for c := range ch {
			hash, err := structhash.Hash(c.Configurations, 1)
			if err != nil {
				log.WithError(err).Error("Could not calculate hash for configuration")
				return
			}

			if s.version == hash {
				log.Debug("Skipping same configuration")
				return
			}

			s.version = hash
			s.handleEvent(c.Configurations)
		}
	}()

	r := s.createRouter()
	// some routers may panic when have empty routes list, so add one dummy 404 route to avoid this
	if r.RoutesCount() < 1 {
		r.Any("/", httpErrors.NotFound)
	}

	s.register = proxy.NewRegister(r, proxy.Params{
		StatsClient:            s.statsClient,
		FlushInterval:          s.globalConfig.BackendFlushInterval,
		IdleConnectionsPerHost: s.globalConfig.MaxIdleConnsPerHost,
		CloseIdleConnsPeriod:   s.globalConfig.CloseIdleConnsPeriod,
	})
	s.defLoader = loader.NewAPILoader(s.register)

	webServer := web.New(
		repo,
		web.WithPort(s.globalConfig.Web.Port),
		web.WithTLS(s.globalConfig.Web.TLS),
		web.WithCredentials(s.globalConfig.Web.Credentials),
		web.ReadOnly(s.globalConfig.Web.ReadOnly),
	)
	if err := webServer.Start(); err != nil {
		return errors.Wrap(err, "could not start Janus web API")
	}

	return s.listenAndServe(r)
}

// Start starts the server
func (s *Server) Start() error {
	return s.StartWithContext(context.Background())
}

// Close closes the server
func (s *Server) Close() error {
	return s.server.Close()
}

func (s *Server) listenAndServe(handler http.Handler) error {
	address := fmt.Sprintf(":%v", s.globalConfig.Port)
	s.server = &http.Server{Addr: address, Handler: handler}

	if s.globalConfig.TLS.IsHTTPS() {
		s.server.Addr = fmt.Sprintf(":%v", s.globalConfig.TLS.Port)

		if s.globalConfig.TLS.Redirect {
			go func() {
				log.WithField("address", address).Info("Listening HTTP redirects to HTTPS")
				log.Fatal(http.ListenAndServe(address, web.RedirectHTTPS(s.globalConfig.TLS.Port)))
			}()
		}

		log.WithField("address", s.server.Addr).Info("Listening HTTPS")
		return s.server.ListenAndServeTLS(s.globalConfig.TLS.CertFile, s.globalConfig.TLS.KeyFile)
	}

	log.WithField("address", address).Info("Certificate and certificate key were not found, defaulting to HTTP")
	return s.server.ListenAndServe()
}

func (s *Server) createRouter() router.Router {
	// create router with a custom not found handler
	router.DefaultOptions.NotFoundHandler = errors.NotFound
	r := router.NewChiRouterWithOptions(router.DefaultOptions)
	r.Use(
		middleware.NewStats(s.statsClient).Handler,
		middleware.NewLogger().Handler,
		middleware.NewRecovery(errors.RecoveryHandler),
		middleware.NewOpenTracing(s.globalConfig.TLS.IsHTTPS()).Handler,
	)

	if s.globalConfig.RequestID {
		r.Use(middleware.RequestID)
	}

	return r
}

func (s *Server) buildConfiguration(defs []*api.Definition) []*api.Spec {
	var specs []*api.Spec
	for _, d := range defs {
		specs = append(specs, &api.Spec{Definition: d})
	}

	return specs
}

func (s *Server) handleEvent(defs []*api.Definition) {
	if !s.started {
		specs := s.buildConfiguration(defs)

		event := plugin.OnStartup{
			StatsClient:   s.statsClient,
			Register:      s.register,
			Config:        s.globalConfig,
			Configuration: specs,
		}

		if mgoRepo, ok := s.repo.(*api.MongoRepository); ok {
			event.MongoSession = mgoRepo.Session
		}

		plugin.EmitEvent(plugin.StartupEvent, event)

		s.defLoader.RegisterAPIs(specs)
		s.started = true
		log.Info("Janus started")
	} else {
		log.Debug("Refreshing configuration")
		newRouter := s.createRouter()
		s.register.UpdateRouter(newRouter)

		specs := s.buildConfiguration(defs)
		s.defLoader.RegisterAPIs(specs)

		plugin.EmitEvent(plugin.ReloadEvent, plugin.OnReload{Configurations: specs})

		s.server.Handler = newRouter
		log.Debug("Configuration refresh done")
	}
}
