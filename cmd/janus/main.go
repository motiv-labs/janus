package main

import (
	"fmt"
	"net/http"

	mgo "gopkg.in/mgo.v2"

	"net/url"

	log "github.com/Sirupsen/logrus"
	"github.com/hellofresh/janus/pkg/api"
	"github.com/hellofresh/janus/pkg/errors"
	"github.com/hellofresh/janus/pkg/middleware"
	"github.com/hellofresh/janus/pkg/oauth"
	"github.com/hellofresh/janus/pkg/proxy"
	"github.com/hellofresh/janus/pkg/router"
	"github.com/hellofresh/janus/pkg/stats"
	"github.com/hellofresh/janus/pkg/web"
)

func main() {
	var repo api.APISpecRepository
	var oAuthServersRepo oauth.Repository
	var err error

	defer statsdClient.Close()

	statsClient := stats.NewStatsClient(statsdClient)
	dsnURL, err := url.Parse(globalConfig.Database.DSN)

	switch dsnURL.Scheme {
	case "mongodb":
		log.Debugf("Trying to connect to %s", globalConfig.Database.DSN)
		session, err := mgo.Dial(globalConfig.Database.DSN)
		defer session.Close()
		if err != nil {
			log.Panic(err)
		}

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
		log.Debug("Loading API definitions from file system")
		repo, err = api.NewFileSystemRepository(dsnURL.Path + "/apis")
		if err != nil {
			log.Panic(err)
		}

		log.Debug("Loading OAuth servers definitions from file system")
		oAuthServersRepo, err = oauth.NewFileSystemRepository(dsnURL.Path + "/auth")
		if err != nil {
			log.Panic(err)
		}
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

	// create router
	r := router.NewHttpTreeMuxRouter()
	r.Use(
		middleware.NewStats(statsClient).Handler,
		middleware.NewLogger().Handler,
		middleware.NewRecovery(web.RecoveryHandler).Handler,
	)

	// create proxy register
	register := proxy.NewRegister(r, p)

	apiLoader := api.NewLoader(register, storage, oAuthServersRepo)
	apiLoader.LoadDefinitions(repo)

	oauthLoader := oauth.NewLoader(register, storage)
	oauthLoader.LoadDefinitions(oAuthServersRepo)

	wp := web.Provider{
		Port:     "3001",
		Cred:     globalConfig.Credentials,
		APIRepo:  repo,
		AuthRepo: oAuthServersRepo,
	}
	wp.Provide()

	log.Fatal(listenAndServe(r))
}

func listenAndServe(handler http.Handler) error {
	address := fmt.Sprintf(":%v", globalConfig.Port)
	log.Infof("Listening on %v", address)
	if globalConfig.IsHTTPS() {
		return http.ListenAndServeTLS(address, globalConfig.CertPathTLS, globalConfig.KeyPathTLS, handler)
	}

	log.Infof("certPathTLS or keyPathTLS not found, defaulting to HTTP")
	return http.ListenAndServe(address, handler)
}
