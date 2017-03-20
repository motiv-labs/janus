package main

import (
	"fmt"
	"net/http"
	"net/url"

	log "github.com/Sirupsen/logrus"
	"github.com/hellofresh/janus/pkg/api"
	"github.com/hellofresh/janus/pkg/errors"
	"github.com/hellofresh/janus/pkg/middleware"
	"github.com/hellofresh/janus/pkg/oauth"
	"github.com/hellofresh/janus/pkg/proxy"
	"github.com/hellofresh/janus/pkg/router"
	"github.com/hellofresh/janus/pkg/web"
	"github.com/hellofresh/stats-go"
	mgo "gopkg.in/mgo.v2"
)

func main() {
	var repo api.Repository
	var oAuthServersRepo oauth.Repository
	var readOnlyAPI bool
	var err error

	log.Info("Janus starting...")

	sectionsTestsMap, err := stats.ParseSectionsTestsMap(globalConfig.Stats.IDs)
	if err != nil {
		log.WithError(err).WithField("config", globalConfig.Stats.IDs).
			Error("Failed to parse stats second level IDs from env")
		sectionsTestsMap = map[stats.PathSection]stats.SectionTestDefinition{}
	}
	log.WithField("config", globalConfig.Stats.IDs).WithField("map", sectionsTestsMap.String()).
		Debug("Setting stats second level IDs")

	statsClient := stats.NewStatsdStatsClient(globalConfig.Stats.DSN, globalConfig.Stats.Prefix).
		SetHttpMetricCallback(stats.NewHasIDAtSecondLevelCallback(sectionsTestsMap))
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

		readOnlyAPI = true
	default:
		log.WithError(errors.ErrInvalidScheme).Error("No Database selected")
	}

	transport := oauth.NewAwareTransport(statsClient, storage, oAuthServersRepo)
	p := proxy.WithParams(proxy.Params{
		Transport:              transport,
		FlushInterval:          globalConfig.BackendFlushInterval,
		IdleConnectionsPerHost: globalConfig.MaxIdleConnsPerHost,
		CloseIdleConnsPeriod:   globalConfig.CloseIdleConnsPeriod,
		InsecureSkipVerify:     globalConfig.InsecureSkipVerify,
	})
	defer p.Close()

	// create router with a custom not found handler
	router.DefaultOptions.NotFoundHandler = web.NotFound
	r := router.NewHttpTreeMuxWithOptions(router.DefaultOptions)
	r.Use(
		middleware.NewStats(statsClient).Handler,
		middleware.NewLogger().Handler,
		middleware.NewRecovery(web.RecoveryHandler).Handler,
		middleware.NewOpenTracing().Handler,
	)

	// create proxy register
	register := proxy.NewRegister(r, p)

	apiLoader := api.NewLoader(register, storage, oAuthServersRepo)
	apiLoader.LoadDefinitions(repo)

	oauthLoader := oauth.NewLoader(register, storage)
	oauthLoader.LoadDefinitions(oAuthServersRepo)

	wp := web.Provider{
		Port:     globalConfig.APIPort,
		Cred:     globalConfig.Credentials,
		APIRepo:  repo,
		AuthRepo: oAuthServersRepo,
		ReadOnly: readOnlyAPI,
	}
	wp.Provide()

	log.Fatal(listenAndServe(r))
}

func listenAndServe(handler http.Handler) error {
	address := fmt.Sprintf(":%v", globalConfig.Port)
	log.WithField("address", address).Info("Listening on")
	log.Info("Janus started")
	if globalConfig.IsHTTPS() {
		return http.ListenAndServeTLS(address, globalConfig.CertPathTLS, globalConfig.KeyPathTLS, handler)
	}

	log.Info("Certificate and certificate key were not found, defaulting to HTTP")
	return http.ListenAndServe(address, handler)
}
