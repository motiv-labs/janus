package main

import (
	"os"
	"path/filepath"

	"github.com/hellofresh/janus/pkg/config"
	tracerfactory "github.com/hellofresh/janus/pkg/opentracing"
	"github.com/hellofresh/janus/pkg/store"
	"github.com/hellofresh/stats-go"
	"github.com/hellofresh/stats-go/bucket"
	"github.com/hellofresh/stats-go/hooks"
	"github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
)

var (
	globalConfig *config.Specification
	statsClient  stats.Client
	storage      store.Store
)

func initConfig() {
	c, err := config.Load(configFile)
	if nil != err {
		log.WithError(err).Panic("Could not parse the environment configurations")
	}

	globalConfig = c
}

// initializes the basic configuration for the log wrapper
func initLog() {
	err := globalConfig.Log.Apply()
	if nil != err {
		log.WithError(err).Panic("Could not apply logging configurations")
	}
}

// initializes distributed tracing
func initDistributedTracing() {
	log.Debug("Initializing distributed tracing")
	tracer, err := tracerfactory.Build(globalConfig.Tracing)
	if err != nil {
		log.WithError(err).Panic("Could not build a tracer")
	}

	opentracing.SetGlobalTracer(tracer)
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
		log.WithError(err).Panic("Error initializing statsd client")
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
}

// initializes the storage and managers
func initStorage() {
	log.WithField("dsn", globalConfig.Storage.DSN).Debug("Initializing storage")
	s, err := store.Build(globalConfig.Storage.DSN)
	if nil != err {
		log.Panic(err)
	}

	storage = s
}
