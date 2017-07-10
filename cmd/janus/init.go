package main

import (
	"github.com/hellofresh/janus/pkg/config"
	tracerfactory "github.com/hellofresh/janus/pkg/opentracing"
	"github.com/hellofresh/janus/pkg/store"
	"github.com/hellofresh/stats-go"
	"github.com/hellofresh/stats-go/bucket"
	"github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
)

var (
	globalConfig *config.Specification
	statsClient  stats.Client
	storage      store.Store
)

func init() {
	c, err := config.Load(configFile)
	if nil != err {
		log.WithError(err).Panic("Could not parse the environment configurations")
	}

	globalConfig = c
}

// initializes the basic configuration for the log wrapper
func init() {
	err := globalConfig.Log.Apply()
	if nil != err {
		log.WithError(err).Panic("Could not apply logging configurations")
	}
}

// initializes distributed tracing
func init() {
	log.Debug("Initializing distributed tracing")
	tracer, err := tracerfactory.Build(globalConfig.Tracing)
	if err != nil {
		log.WithError(err).Panic("Could not build a tracer")
	}

	opentracing.SetGlobalTracer(tracer)
}

func init() {
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

	statsClient.SetHTTPMetricCallback(bucket.NewHasIDAtSecondLevelCallback(sectionsTestsMap))
}

// initializes the storage and managers
func init() {
	log.WithField("dsn", globalConfig.Storage.DSN).Debug("Initializing storage")
	s, err := store.Build(globalConfig.Storage.DSN)
	if nil != err {
		log.Panic(err)
	}

	storage = s
}
