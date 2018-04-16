package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"reflect"
	"time"

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
	server                *http.Server
	provider              api.Repository
	register              *proxy.Register
	defLoader             *loader.APILoader
	currentConfigurations *api.Configuration
	configurationChan     chan api.ConfigurationChanged
	stopChan              chan struct{}
	globalConfig          *config.Specification
	statsClient           client.Client
	webServer             *web.Server
}

// New creates a new instance of Server
func New(opts ...Option) *Server {
	s := Server{
		configurationChan: make(chan api.ConfigurationChanged, 100),
		stopChan:          make(chan struct{}, 1),
	}

	for _, opt := range opts {
		opt(&s)
	}

	return &s
}

// Start starts the server
func (s *Server) Start() error {
	return s.StartWithContext(context.Background())
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
	}()
	go func() {
		if err := s.startHTTPServers(); err != nil {
			log.WithError(err).Fatal("Could not start http servers")
		}
	}()

	go s.listenProviders(s.stopChan)

	defs, err := s.provider.FindAll()
	if err != nil {
		return errors.Wrap(err, "could not find all configurations from the provider")
	}

	s.currentConfigurations = &api.Configuration{Definitions: defs}
	if err := s.startProvider(ctx); err != nil {
		log.WithError(err).Fatal("Could not start providers")
	}

	event := plugin.OnStartup{
		StatsClient:   s.statsClient,
		Register:      s.register,
		Config:        s.globalConfig,
		Configuration: defs,
	}

	if mgoRepo, ok := s.provider.(*api.MongoRepository); ok {
		event.MongoSession = mgoRepo.Session
	}

	plugin.EmitEvent(plugin.StartupEvent, event)
	s.defLoader.RegisterAPIs(defs)

	log.Info("Janus started")

	return nil
}

// Wait blocks until server is shutted down.
func (s *Server) Wait() {
	<-s.stopChan
}

// Stop stops the server
func (s *Server) Stop() {
	defer log.Info("Server stopped")

	graceTimeOut := time.Duration(s.globalConfig.GraceTimeOut)
	ctx, cancel := context.WithTimeout(context.Background(), graceTimeOut)
	defer cancel()
	log.Debugf("Waiting %s seconds before killing connections...", graceTimeOut)
	if err := s.server.Shutdown(ctx); err != nil {
		log.WithError(err).Debug("Wait is over due to error")
		s.server.Close()
	}
	log.Debug("Server closed")

	s.stopChan <- struct{}{}
}

// Close closes the server
func (s *Server) Close() error {
	defer close(s.stopChan)
	defer close(s.configurationChan)
	defer s.webServer.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	go func(ctx context.Context) {
		<-ctx.Done()
		if ctx.Err() == context.Canceled {
			return
		} else if ctx.Err() == context.DeadlineExceeded {
			panic("Timeout while stopping janus, killing instance ✝")
		}
	}(ctx)

	return s.server.Close()
}

func (s *Server) startHTTPServers() error {
	r := s.createRouter()
	s.register = proxy.NewRegister(r, proxy.Params{
		StatsClient:            s.statsClient,
		FlushInterval:          s.globalConfig.BackendFlushInterval,
		IdleConnectionsPerHost: s.globalConfig.MaxIdleConnsPerHost,
		CloseIdleConnsPeriod:   s.globalConfig.CloseIdleConnsPeriod,
	})
	s.defLoader = loader.NewAPILoader(s.register)

	return s.listenAndServe(r)
}

func (s *Server) startProvider(ctx context.Context) error {
	s.webServer = web.New(
		web.WithConfigurations(s.currentConfigurations),
		web.WithPort(s.globalConfig.Web.Port),
		web.WithTLS(s.globalConfig.Web.TLS),
		web.WithCredentials(s.globalConfig.Web.Credentials),
		web.ReadOnly(s.globalConfig.Web.ReadOnly),
	)

	if err := s.webServer.Start(); err != nil {
		return errors.Wrap(err, "could not start Janus web API")
	}

	if listener, ok := s.provider.(api.Listener); ok {
		go func() {
			ch := make(chan api.ConfigurationMessage)
			listener.Listen(ctx, ch)
			for {
				select {
				case c, more := <-s.webServer.ConfigurationChan:
					if !more {
						return
					}

					s.updateConfigurations(c)
					s.handleEvent(s.currentConfigurations)
					ch <- c
				case <-ctx.Done():
					close(ch)
					return
				}
			}
		}()
	}

	if watcher, ok := s.provider.(api.Watcher); ok {
		watcher.Watch(ctx, s.configurationChan)
	}

	return nil
}

func (s *Server) listenProviders(stop chan struct{}) {
	for {
		select {
		case <-stop:
			return
		case configMsg, ok := <-s.configurationChan:
			if !ok {
				return
			}

			if reflect.DeepEqual(s.currentConfigurations, configMsg.Configurations) {
				log.Debug("Skipping same configuration")
				continue
			}

			s.currentConfigurations.Definitions = configMsg.Configurations.Definitions
			s.handleEvent(configMsg.Configurations)
		}
	}
}

func (s *Server) listenAndServe(handler http.Handler) error {
	address := fmt.Sprintf(":%v", s.globalConfig.Port)
	logger := log.WithField("address", address)
	s.server = &http.Server{Addr: address, Handler: handler}
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return errors.Wrap(err, "error opening listener")
	}

	if s.globalConfig.TLS.IsHTTPS() {
		s.server.Addr = fmt.Sprintf(":%v", s.globalConfig.TLS.Port)

		if s.globalConfig.TLS.Redirect {
			go func() {
				logger.Info("Listening HTTP redirects to HTTPS")
				log.Fatal(http.Serve(listener, web.RedirectHTTPS(s.globalConfig.TLS.Port)))
			}()
		}

		logger.Info("Listening HTTPS")
		return s.server.ServeTLS(listener, s.globalConfig.TLS.CertFile, s.globalConfig.TLS.KeyFile)
	}

	logger.Info("Certificate and certificate key were not found, defaulting to HTTP")
	return s.server.Serve(listener)
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

	// some routers may panic when have empty routes list, so add one dummy 404 route to avoid this
	if r.RoutesCount() < 1 {
		r.Any("/", httpErrors.NotFound)
	}

	return r
}

func (s *Server) updateConfigurations(cfg api.ConfigurationMessage) {
	currentDefinitions := s.currentConfigurations.Definitions

	switch cfg.Operation {
	case api.AddedOperation:
		currentDefinitions = append(currentDefinitions, cfg.Configuration)
	case api.UpdatedOperation:
		for i, d := range currentDefinitions {
			if d.Name == cfg.Configuration.Name {
				currentDefinitions[i] = cfg.Configuration
			}
		}
	case api.RemovedOperation:
		for i, d := range currentDefinitions {
			if d.Name == cfg.Configuration.Name {
				copy(currentDefinitions[i:], currentDefinitions[i+1:])
				// currentDefinitions[len(currentDefinitions)-1] = nil // or the zero value of T
				currentDefinitions = currentDefinitions[:len(currentDefinitions)-1]
			}
		}
	}

	s.currentConfigurations.Definitions = currentDefinitions
}

func (s *Server) handleEvent(cfg *api.Configuration) {
	log.Debug("Refreshing configuration")
	newRouter := s.createRouter()

	s.register.UpdateRouter(newRouter)
	s.defLoader.RegisterAPIs(cfg.Definitions)

	plugin.EmitEvent(plugin.ReloadEvent, plugin.OnReload{Configurations: cfg.Definitions})

	s.server.Handler = newRouter
	log.Debug("Configuration refresh done")
}
