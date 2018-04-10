package main

import (
	"fmt"
	"net/http"

	"github.com/hellofresh/janus/pkg/api"
	"github.com/hellofresh/janus/pkg/cluster"
	"github.com/hellofresh/janus/pkg/errors"
	"github.com/hellofresh/janus/pkg/middleware"
	"github.com/hellofresh/janus/pkg/opentracing"
	"github.com/hellofresh/janus/pkg/plugin"
	"github.com/hellofresh/janus/pkg/proxy"
	"github.com/hellofresh/janus/pkg/router"
	"github.com/hellofresh/janus/pkg/web"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	// this is needed to call the init function on each plugin
	_ "github.com/hellofresh/janus/pkg/plugin/basic"
	_ "github.com/hellofresh/janus/pkg/plugin/bodylmt"
	_ "github.com/hellofresh/janus/pkg/plugin/compression"
	_ "github.com/hellofresh/janus/pkg/plugin/cors"
	_ "github.com/hellofresh/janus/pkg/plugin/oauth2"
	_ "github.com/hellofresh/janus/pkg/plugin/rate"
	_ "github.com/hellofresh/janus/pkg/plugin/requesttransformer"
	_ "github.com/hellofresh/janus/pkg/plugin/responsetransformer"

	// dynamically registered auth providers
	_ "github.com/hellofresh/janus/pkg/jwt/basic"
	_ "github.com/hellofresh/janus/pkg/jwt/github"

	// internal plugins
	_ "github.com/hellofresh/janus/pkg/loader"
	_ "github.com/hellofresh/janus/pkg/web"
)

var (
	server *http.Server
)

// RunServer is the run command to start Janus
func RunServer(cmd *cobra.Command, args []string) {
	log.WithField("version", version).Info("Janus starting...")

	initConfig()
	initLog()
	initStatsd()
	initDatabase()

	tracingFactory := opentracing.New(globalConfig.Tracing)
	tracingFactory.Setup()

	defer tracingFactory.Close()
	defer statsClient.Close()
	defer globalConfig.Log.Flush()
	defer session.Close()

	repo, err := api.BuildRepository(globalConfig.Database.DSN, session)
	if err != nil {
		log.Panic(err)
	}

	cluster.Watch(globalConfig.Cluster.UpdateFrequency, handleEvent(repo))

	r := createRouter()
	register := proxy.NewRegister(r, proxy.Params{
		StatsClient:            statsClient,
		FlushInterval:          globalConfig.BackendFlushInterval,
		IdleConnectionsPerHost: globalConfig.MaxIdleConnsPerHost,
		CloseIdleConnsPeriod:   globalConfig.CloseIdleConnsPeriod,
	})

	webServer := web.New(
		repo,
		web.WithPort(globalConfig.Web.Port),
		web.WithTLS(globalConfig.Web.TLS),
		web.WithCredentials(globalConfig.Web.Credentials),
		web.ReadOnly(globalConfig.Web.ReadOnly),
	)
	if err := webServer.Serve(); err != nil {
		log.WithError(err).Fatal("Could not start Janus web API")
	}

	configuration := buildConfiguration(repo)

	event := plugin.OnStartup{
		StatsClient:   statsClient,
		MongoSession:  session,
		Register:      register,
		Config:        globalConfig,
		Configuration: configuration,
	}
	plugin.EmitEvent(plugin.StartupEvent, event)

	log.Fatal(listenAndServe(r))
}

func buildConfiguration(repo api.Repository) []*api.Spec {
	defs, err := repo.FindAll()
	if err != nil {
		log.Panic(err)
	}

	var specs []*api.Spec
	for _, definition := range defs {
		specs = append(specs, &api.Spec{Definition: definition})
	}

	return specs
}

func listenAndServe(handler http.Handler) error {
	address := fmt.Sprintf(":%v", globalConfig.Port)
	server = &http.Server{Addr: address, Handler: handler}

	log.Info("Janus started")
	if globalConfig.TLS.IsHTTPS() {
		server.Addr = fmt.Sprintf(":%v", globalConfig.TLS.Port)

		if globalConfig.TLS.Redirect {
			go func() {
				log.WithField("address", address).Info("Listening HTTP redirects to HTTPS")
				log.Fatal(http.ListenAndServe(address, web.RedirectHTTPS(globalConfig.TLS.Port)))
			}()
		}

		log.WithField("address", server.Addr).Info("Listening HTTPS")
		return server.ListenAndServeTLS(globalConfig.TLS.CertFile, globalConfig.TLS.KeyFile)
	}

	log.WithField("address", address).Info("Certificate and certificate key were not found, defaulting to HTTP")
	return server.ListenAndServe()
}

func createRouter() router.Router {
	// create router with a custom not found handler
	router.DefaultOptions.NotFoundHandler = errors.NotFound
	r := router.NewChiRouterWithOptions(router.DefaultOptions)
	r.Use(
		middleware.NewStats(statsClient).Handler,
		middleware.NewLogger().Handler,
		middleware.NewRecovery(errors.RecoveryHandler),
		middleware.NewOpenTracing(globalConfig.TLS.IsHTTPS()).Handler,
	)

	if globalConfig.RequestID {
		r.Use(middleware.RequestID)
	}

	return r
}

func handleEvent(repo api.Repository) func() {
	return func() {
		log.Debug("Refreshing configuration")
		specs := buildConfiguration(repo)

		event := plugin.OnReload{Configurations: specs}

		plugin.EmitEvent(plugin.ReloadEvent, event)
		log.Debug("Configuration refresh done")
	}
}
