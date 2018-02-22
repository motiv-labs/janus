package main

import (
	"io"
	"net/url"
	"os"
	"path/filepath"

	"github.com/afex/hystrix-go/hystrix/metric_collector"
	"github.com/afex/hystrix-go/plugins"
	"github.com/hellofresh/janus/pkg/config"
	tracerfactory "github.com/hellofresh/janus/pkg/opentracing"
	"github.com/hellofresh/janus/pkg/store"
	"github.com/hellofresh/stats-go"
	"github.com/hellofresh/stats-go/bucket"
	"github.com/hellofresh/stats-go/hooks"
	opentracing "github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
)

var (
	globalConfig *config.Specification
	statsClient  stats.Client
	storage      store.Store
	session      *mgo.Session
)

func initConfig() {
	var err error
	globalConfig, err = config.Load(configFile)
	if nil != err {
		log.WithError(err).Error("Could not load configurations from file - trying environment configurations instead.")

		globalConfig, err = config.LoadEnv()
		if nil != err {
			log.WithError(err).Error("Could not load configurations from environment")
		}
	}
}

// initializes the basic configuration for the log wrapper
func initLog() {
	err := globalConfig.Log.Apply()
	if nil != err {
		log.WithError(err).Fatal("Could not apply logging configurations")
	}
}

// initializes distributed tracing
func initDistributedTracing() io.Closer {
	log.Debug("Initializing distributed tracing")
	tracer, closer, err := tracerfactory.Build(globalConfig.Tracing)
	if err != nil {
		log.WithError(err).Fatal("Could not build a tracer")
	}

	opentracing.SetGlobalTracer(tracer)

	return closer
}

func initStatsd() {
	sectionsTestsMap, err := bucket.ParseSectionsTestsMap(globalConfig.Stats.IDs)
	if err != nil {
		log.WithError(err).WithField("config", globalConfig.Stats.IDs).
			Error("Failed to parse stats second level IDs from env")
		sectionsTestsMap = map[bucket.PathSection]bucket.SectionTestDefinition{}
	}
	log.WithField("config", globalConfig.Stats.IDs).
		WithField("map", sectionsTestsMap.String()).
		Debug("Setting stats second level IDs")

	statsClient, err = stats.NewClient(globalConfig.Stats.DSN, globalConfig.Stats.Prefix)
	if err != nil {
		log.WithError(err).Fatal("Error initializing statsd client")
	}

	statsClient.SetHTTPMetricCallback(bucket.NewHasIDAtSecondLevelCallback(&bucket.SecondLevelIDConfig{
		HasIDAtSecondLevel:    sectionsTestsMap,
		AutoDiscoverThreshold: globalConfig.Stats.AutoDiscoverThreshold,
		AutoDiscoverWhiteList: globalConfig.Stats.AutoDiscoverWhiteList,
	}))

	host, err := os.Hostname()
	if nil != err {
		host = "-unknown-"
	}

	_, appFile := filepath.Split(os.Args[0])
	statsClient.TrackMetric("app", bucket.MetricOperation{"init", host, appFile})
	log.AddHook(hooks.NewLogrusHook(statsClient, globalConfig.Stats.ErrorsSection))

	// Setup metrics for circuit breaker
	c, err := plugins.InitializeStatsdCollector(&plugins.StatsdCollectorConfig{
		StatsdAddr: globalConfig.Stats.DSN,
		Prefix:     globalConfig.Stats.Prefix,
	})
	metricCollector.Registry.Register(c.NewStatsdCollector)
	if err != nil {
		log.WithError(err).Error("Error initializing statsd client for circuit breaker")
	}
}

// initializes the storage and managers
func initStorage() {
	log.WithField("dsn", globalConfig.Storage.DSN).Debug("Initializing storage")
	s, err := store.Build(globalConfig.Storage.DSN)
	if nil != err {
		log.Fatal(err)
	}

	storage = s
}

// initializes the storage and managers
func initDatabase() {
	dsnURL, err := url.Parse(globalConfig.Database.DSN)
	switch dsnURL.Scheme {
	case "mongodb":
		log.Debug("MongoDB configuration chosen")

		log.WithField("dsn", globalConfig.Database.DSN).Debug("Trying to connect to MongoDB...")
		session, err = mgo.Dial(globalConfig.Database.DSN)
		if err != nil {
			log.Fatal(err)
		}

		log.Debug("Connected to MongoDB")
		session.SetMode(mgo.Monotonic, true)
	default:
		log.Error("No Database selected")
	}
}
