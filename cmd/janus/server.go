package main

import (
	"fmt"
	"net/http"

	"github.com/hellofresh/janus/pkg/api"
	"github.com/hellofresh/janus/pkg/errors"
	"github.com/hellofresh/janus/pkg/middleware"
	"github.com/hellofresh/janus/pkg/notifier"
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
	var ntf notifier.Notifier

	log.WithField("version", version).Info("Janus starting...")

	initConfig()
	initLog()
	initStatsd()
	initStorage()
	initDatabase()
	dtCloser := initDistributedTracing()

	defer dtCloser.Close()
	defer statsClient.Close()
	defer globalConfig.Log.Flush()
	defer session.Close()

	repo, err := api.BuildRepository(globalConfig.Database.DSN, session)
	if err != nil {
		log.Panic(err)
	}

	if subscriber, ok := storage.(notifier.Subscriber); ok {
		listener := notifier.NewNotificationListener(subscriber)
		listener.Start(handleEvent(repo))
	}

	if publisher, ok := storage.(notifier.Publisher); ok {
		ntf = notifier.NewPublisherNotifier(publisher, "")
	}

	r := createRouter()
	register := proxy.NewRegister(r, proxy.Params{
		StatsClient:            statsClient,
		FlushInterval:          globalConfig.BackendFlushInterval,
		IdleConnectionsPerHost: globalConfig.MaxIdleConnsPerHost,
		CloseIdleConnsPeriod:   globalConfig.CloseIdleConnsPeriod,
	})

	event := plugin.OnStartup{
		Notifier:     ntf,
		Repository:   repo,
		StatsClient:  statsClient,
		MongoSession: session,
		Register:     register,
		Config:       globalConfig,
	}
	plugin.EmitEvent(plugin.StartupEvent, event)

	log.Fatal(listenAndServe(r))
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
	return r
}

func handleEvent(repo api.Repository) func(notification notifier.Notification) {
	return func(notification notifier.Notification) {
		if notifier.RequireReload(notification.Command) {
			newRouter := createRouter()
			register := proxy.NewRegister(newRouter, proxy.Params{
				StatsClient:            statsClient,
				FlushInterval:          globalConfig.BackendFlushInterval,
				IdleConnectionsPerHost: globalConfig.MaxIdleConnsPerHost,
				CloseIdleConnsPeriod:   globalConfig.CloseIdleConnsPeriod,
			})

			event := plugin.OnReload{
				Register:   register,
				Repository: repo,
			}

			plugin.EmitEvent(plugin.ReloadEvent, event)

			server.Handler = newRouter
		}
	}
}
