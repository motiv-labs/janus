package main

import (
	"fmt"
	"net/http"
	"net/url"

	log "github.com/Sirupsen/logrus"
	"github.com/hellofresh/janus/pkg/api"
	"github.com/hellofresh/janus/pkg/errors"
	"github.com/hellofresh/janus/pkg/loader"
	"github.com/hellofresh/janus/pkg/middleware"
	"github.com/hellofresh/janus/pkg/notifier"
	"github.com/hellofresh/janus/pkg/oauth"
	"github.com/hellofresh/janus/pkg/proxy"
	"github.com/hellofresh/janus/pkg/router"
	"github.com/hellofresh/janus/pkg/web"
	"github.com/spf13/cobra"
	mgo "gopkg.in/mgo.v2"
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

	if subscriber, ok := storage.(notifier.Subscriber); ok {
		listerner := notifier.NewNotificationListener(subscriber)
		listerner.Start(handleEvent)
	}

	dsnURL, err := url.Parse(globalConfig.Database.DSN)
	switch dsnURL.Scheme {
	case "mongodb":
		log.WithField("dsn", globalConfig.Database.DSN).Debug("Trying to connect to DB")
		session, err := mgo.Dial(globalConfig.Database.DSN)
		if err != nil {
			log.Panic(err)
		}

		defer session.Close()

		log.Debug("Connected to mongodb")
		session.SetMode(mgo.Monotonic, true)

		log.Debug("Loading API definitions from Mongo DB")
		repo, err = api.NewMongoAppRepository(session)
		if err != nil {
			log.Panic(err)
		}

		// create the proxy
		log.Debug("Loading OAuth servers definitions from Mongo DB")
		oAuthServersRepo, err = oauth.NewMongoRepository(session)
		if err != nil {
			log.Panic(err)
		}
	case "file":
		var apiPath = dsnURL.Path + "/apis"
		var authPath = dsnURL.Path + "/auth"

		log.WithField("path", apiPath).Debug("Loading API definitions from file system")
		repo, err = api.NewFileSystemRepository(apiPath)
		if err != nil {
			log.Panic(err)
		}

		log.WithField("path", authPath).Debug("Loading OAuth servers definitions from file system")
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
		CertFile: globalConfig.Web.CertFile,
		KeyFile:  globalConfig.Web.KeyFile,
		APIRepo:  repo,
		AuthRepo: oAuthServersRepo,
		ReadOnly: globalConfig.Web.ReadOnly,
	}

	if publisher, ok := storage.(notifier.Publisher); ok {
		wp.Notifier = notifier.New(publisher, "")
	}

	wp.Provide(version)

	r := createRouter()

	loader.Load(loader.Params{
		Router:    r,
		Storage:   storage,
		APIRepo:   repo,
		OAuthRepo: oAuthServersRepo,
		ProxyParams: proxy.Params{
			StatsClient:            statsClient,
			FlushInterval:          globalConfig.BackendFlushInterval,
			IdleConnectionsPerHost: globalConfig.MaxIdleConnsPerHost,
			CloseIdleConnsPeriod:   globalConfig.CloseIdleConnsPeriod,
			InsecureSkipVerify:     globalConfig.InsecureSkipVerify,
		},
	})

	address := fmt.Sprintf(":%v", globalConfig.Port)
	log.WithField("address", address).Info("Listening on")
	server = &http.Server{Addr: address, Handler: r}
	log.Fatal(listenAndServe(server))
}

func listenAndServe(server *http.Server) error {
	log.Info("Janus started")
	if globalConfig.Web.IsHTTPS() {
		return server.ListenAndServeTLS(globalConfig.CertFile, globalConfig.KeyFile)
	}

	log.Info("Certificate and certificate key were not found, defaulting to HTTP")
	return server.ListenAndServe()
}

func createRouter() router.Router {
	// create router with a custom not found handler
	router.DefaultOptions.NotFoundHandler = web.NotFound
	r := router.NewChiRouterWithOptions(router.DefaultOptions)
	r.Use(
		middleware.NewStats(statsClient).Handler,
		middleware.NewLogger().Handler,
		middleware.NewRecovery(web.RecoveryHandler).Handler,
		middleware.NewOpenTracing(globalConfig.Web.IsHTTPS()).Handler,
	)
	return r
}

func handleEvent(notif notifier.Notification) {
	if notifier.RequireReload(notif.Command) {
		newRouter := createRouter()
		loader.Load(loader.Params{
			Router:    newRouter,
			Storage:   storage,
			APIRepo:   repo,
			OAuthRepo: oAuthServersRepo,
			ProxyParams: proxy.Params{
				StatsClient:            statsClient,
				FlushInterval:          globalConfig.BackendFlushInterval,
				IdleConnectionsPerHost: globalConfig.MaxIdleConnsPerHost,
				CloseIdleConnsPeriod:   globalConfig.CloseIdleConnsPeriod,
				InsecureSkipVerify:     globalConfig.InsecureSkipVerify,
			},
		})
		server.Handler = newRouter
	}
}
