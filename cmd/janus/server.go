package main

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/hellofresh/janus/pkg/api"
	"github.com/hellofresh/janus/pkg/errors"
	"github.com/hellofresh/janus/pkg/loader"
	"github.com/hellofresh/janus/pkg/middleware"
	"github.com/hellofresh/janus/pkg/notifier"
	"github.com/hellofresh/janus/pkg/oauth"
	"github.com/hellofresh/janus/pkg/proxy"
	"github.com/hellofresh/janus/pkg/router"
	"github.com/hellofresh/janus/pkg/web"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/mgo.v2"
)

var (
	repo             api.Repository
	oAuthServersRepo oauth.Repository
	server           *http.Server
)

// RunServer is the run command to start Janus
func RunServer(cmd *cobra.Command, args []string) {
	log.WithField("version", version).Info("Janus starting...")

	defer statsClient.Close()
	defer globalConfig.Log.Flush()

	if subscriber, ok := storage.(notifier.Subscriber); ok {
		listener := notifier.NewNotificationListener(subscriber)
		listener.Start(handleEvent)
	}

	dsnURL, err := url.Parse(globalConfig.Database.DSN)
	switch dsnURL.Scheme {
	case "mongodb":
		log.Debug("MongoDB configuration chosen")

		log.WithField("dsn", globalConfig.Database.DSN).Debug("Trying to connect to MongoDB...")
		session, err := mgo.Dial(globalConfig.Database.DSN)
		if err != nil {
			log.Panic(err)
		}
		defer session.Close()

		log.Debug("Connected to MongoDB")
		session.SetMode(mgo.Monotonic, true)

		repo, err = api.NewMongoAppRepository(session)
		if err != nil {
			log.Panic(err)
		}

		oAuthServersRepo, err = oauth.NewMongoRepository(session)
		if err != nil {
			log.Panic(err)
		}
	case "file":
		log.Debug("File system based configuration chosen")
		var apiPath = dsnURL.Path + "/apis"
		var authPath = dsnURL.Path + "/auth"

		log.WithFields(log.Fields{
			"api_path":  apiPath,
			"auth_path": authPath,
		}).Debug("Trying to load configuration files")
		repo, err = api.NewFileSystemRepository(apiPath)
		if err != nil {
			log.Panic(err)
		}

		oAuthServersRepo, err = oauth.NewFileSystemRepository(authPath)
		if err != nil {
			log.Panic(err)
		}
	default:
		log.WithError(errors.ErrInvalidScheme).Error("No Database selected")
	}

	wp := web.Provider{
		Port:     globalConfig.Web.Port,
		Cred:     globalConfig.Web.Credentials,
		ReadOnly: globalConfig.Web.ReadOnly,
		TLS:      globalConfig.Web.TLS,
		APIRepo:  repo,
		AuthRepo: oAuthServersRepo,
	}

	if publisher, ok := storage.(notifier.Publisher); ok {
		wp.Notifier = notifier.NewPublisherNotifier(publisher, "")
	}

	wp.Provide(version)

	r := createRouter()
	loader.Load(loader.Params{
		Router:      r,
		Storage:     storage,
		APIRepo:     repo,
		OAuthRepo:   oAuthServersRepo,
		StatsClient: statsClient,
		ProxyParams: proxy.Params{
			StatsClient:            statsClient,
			FlushInterval:          globalConfig.BackendFlushInterval,
			IdleConnectionsPerHost: globalConfig.MaxIdleConnsPerHost,
			CloseIdleConnsPeriod:   globalConfig.CloseIdleConnsPeriod,
		},
	})

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

func handleEvent(notification notifier.Notification) {
	if notifier.RequireReload(notification.Command) {
		newRouter := createRouter()
		loader.Load(loader.Params{
			Router:      newRouter,
			Storage:     storage,
			APIRepo:     repo,
			OAuthRepo:   oAuthServersRepo,
			StatsClient: statsClient,
			ProxyParams: proxy.Params{
				StatsClient:            statsClient,
				FlushInterval:          globalConfig.BackendFlushInterval,
				IdleConnectionsPerHost: globalConfig.MaxIdleConnsPerHost,
				CloseIdleConnsPeriod:   globalConfig.CloseIdleConnsPeriod,
			},
		})
		server.Handler = newRouter
	}
}
