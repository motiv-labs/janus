package main

import (
	"fmt"
	"net/http"
	"net/url"

	log "github.com/Sirupsen/logrus"
	"github.com/hellofresh/janus/pkg/api"
	"github.com/hellofresh/janus/pkg/config"
	"github.com/hellofresh/janus/pkg/errors"
	"github.com/hellofresh/janus/pkg/loader"
	"github.com/hellofresh/janus/pkg/middleware"
	"github.com/hellofresh/janus/pkg/oauth"
	"github.com/hellofresh/janus/pkg/plugin"
	"github.com/hellofresh/janus/pkg/proxy"
	"github.com/hellofresh/janus/pkg/router"
	"github.com/hellofresh/janus/pkg/store"
	"github.com/hellofresh/janus/pkg/web"
	"github.com/spf13/cobra"
	mgo "gopkg.in/mgo.v2"
)

var (
	repo             api.Repository
	oAuthServersRepo oauth.Repository
	err              error
	globalConfig     *config.Specification
)

// RunServer is the run command to start Janus
func RunServer(cmd *cobra.Command, args []string) {
	log.Info("Janus starting...")

	globalConfig, err = config.Load(configFile)
	if nil != err {
		log.WithError(err).Panic("Could not parse the environment configurations")
	}

	initLogger(globalConfig.LogLevel)
	initTracing(globalConfig.Tracing)
	storage := initStorage(globalConfig.Storage)
	statsClient := initStatsdClient(globalConfig.Stats)
	defer statsClient.Close()

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

	p := proxy.WithParams(proxy.Params{
		StatsClient:            statsClient,
		FlushInterval:          globalConfig.BackendFlushInterval,
		IdleConnectionsPerHost: globalConfig.MaxIdleConnsPerHost,
		CloseIdleConnsPeriod:   globalConfig.CloseIdleConnsPeriod,
		InsecureSkipVerify:     globalConfig.InsecureSkipVerify,
	})
	defer p.Close()

	// create router with a custom not found handler
	router.DefaultOptions.NotFoundHandler = web.NotFound
	r := router.NewHTTPTreeMuxWithOptions(router.DefaultOptions)
	r.Use(
		middleware.NewStats(statsClient).Handler,
		middleware.NewLogger().Handler,
		middleware.NewRecovery(web.RecoveryHandler).Handler,
		middleware.NewOpenTracing(globalConfig.Web.IsHTTPS()).Handler,
	)

	pluginLoader := plugin.NewLoader()
	pluginLoader.Add(
		plugin.NewRateLimit(storage),
		plugin.NewCORS(),
		plugin.NewOAuth2(oAuthServersRepo, storage),
		plugin.NewCompression(),
	)

	// create proxy register
	register := proxy.NewRegister(r, p)
	var (
		apiSubs   *store.Subscription
		oauthSubs *store.Subscription
	)

	if subscriber, ok := storage.(store.Subscriber); ok {
		apiSubs = subscriber.Subscribe("api_updates")
		oauthSubs = subscriber.Subscribe("oauth_updates")
	}

	apiLoader := loader.NewAPILoader(register, pluginLoader, apiSubs)
	apiLoader.LoadDefinitions(repo)

	oauthLoader := loader.NewOAuthLoader(register, storage, oauthSubs)
	oauthLoader.LoadDefinitions(oAuthServersRepo)

	wp := web.Provider{
		Port:     globalConfig.Web.Port,
		Cred:     globalConfig.Web.Credentials,
		CertFile: globalConfig.Web.CertFile,
		KeyFile:  globalConfig.Web.KeyFile,
		APIRepo:  repo,
		AuthRepo: oAuthServersRepo,
		ReadOnly: globalConfig.Web.ReadOnly,
	}

	if publisher, ok := storage.(store.Publisher); ok {
		wp.Publisher = publisher
	}

	wp.Provide()

	log.Fatal(listenAndServe(r))
}

func listenAndServe(handler http.Handler) error {
	address := fmt.Sprintf(":%v", globalConfig.Port)
	log.WithField("address", address).Info("Listening on")
	log.Info("Janus started")
	if globalConfig.Web.IsHTTPS() {
		return http.ListenAndServeTLS(address, globalConfig.CertFile, globalConfig.KeyFile, handler)
	}

	log.Info("Certificate and certificate key were not found, defaulting to HTTP")
	return http.ListenAndServe(address, handler)
}
