package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/garyburd/redigo/redis"
	"github.com/hellofresh/janus/pkg/api"
	"github.com/hellofresh/janus/pkg/config"
	"github.com/hellofresh/janus/pkg/errors"
	"github.com/hellofresh/janus/pkg/loader"
	"github.com/hellofresh/janus/pkg/middleware"
	"github.com/hellofresh/janus/pkg/notifier"
	"github.com/hellofresh/janus/pkg/oauth"
	"github.com/hellofresh/janus/pkg/plugin"
	"github.com/hellofresh/janus/pkg/proxy"
	"github.com/hellofresh/janus/pkg/router"
	"github.com/hellofresh/janus/pkg/store"
	"github.com/hellofresh/janus/pkg/web"
	stats "github.com/hellofresh/stats-go"
	"github.com/spf13/cobra"
	mgo "gopkg.in/mgo.v2"
)

const (
	PubSubChannel = "janus.cluster.notifications"
)

var (
	repo             api.Repository
	oAuthServersRepo oauth.Repository
	err              error
	globalConfig     *config.Specification
	statsClient      stats.Client
	storage          store.Store
	server           *http.Server
	prx              *proxy.Proxy
	pluginLoader     *plugin.Loader
)

// RunServer is the run command to start Janus
func RunServer(cmd *cobra.Command, args []string) {
	log.WithField("version", version).Info("Janus starting...")

	globalConfig, err = config.Load(configFile)
	if nil != err {
		log.WithError(err).Panic("Could not parse the environment configurations")
	}

	initLogger(globalConfig.LogLevel)
	initTracing(globalConfig.Tracing)
	storage = initStorage(globalConfig.Storage)
	statsClient = initStatsdClient(globalConfig.Stats)
	defer statsClient.Close()

	go startPubSubLoop(storage)

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

	prx = proxy.WithParams(proxy.Params{
		StatsClient:            statsClient,
		FlushInterval:          globalConfig.BackendFlushInterval,
		IdleConnectionsPerHost: globalConfig.MaxIdleConnsPerHost,
		CloseIdleConnsPeriod:   globalConfig.CloseIdleConnsPeriod,
		InsecureSkipVerify:     globalConfig.InsecureSkipVerify,
	})
	defer prx.Close()

	pluginLoader = plugin.NewLoader()
	pluginLoader.Add(
		plugin.NewRateLimit(storage),
		plugin.NewCORS(),
		plugin.NewOAuth2(oAuthServersRepo, storage),
		plugin.NewCompression(),
	)

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
		wp.Notifier = notifier.New(publisher, PubSubChannel)
	}

	wp.Provide(version)

	r := load()
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

func load() router.Router {
	r := createRouter()

	// create proxy register
	register := proxy.NewRegister(r, prx)

	apiLoader := loader.NewAPILoader(register, pluginLoader)
	apiLoader.LoadDefinitions(repo)

	oauthLoader := loader.NewOAuthLoader(register, storage)
	oauthLoader.LoadDefinitions(oAuthServersRepo)

	return r
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

func startPubSubLoop(storage store.Store) {
	log.Debug("Listening for changes")

	if subscriber, ok := storage.(store.Subscriber); ok {
		for {
			err := subscriber.Subscribe(PubSubChannel, func(v interface{}) {
				handleEvent(v)
			})
			if err != nil {
				log.WithFields(log.Fields{
					"prefix": "pub-sub",
					"err":    err,
				}).Error("Connection failed, reconnect in 10s")

				time.Sleep(10 * time.Second)
				log.WithFields(log.Fields{
					"prefix": "pub-sub",
				}).Warning("Reconnecting")
			}
		}
	}
}

func handleEvent(v interface{}) {
	message, ok := v.(redis.Message)
	if !ok {
		return
	}

	notif := notifier.Notification{}
	if err := json.Unmarshal(message.Data, &notif); err != nil {
		log.Error("Unmarshalling message body failed, malformed: ", err)
		return
	}

	if notifier.RequireReload(notif.Command) {
		newRouter := load()
		server.Handler = newRouter
	}
}
