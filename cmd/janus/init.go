package main

import (
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/bshuster-repo/logrus-logstash-hook"
	"github.com/hellofresh/janus/pkg/config"
	tracerfactory "github.com/hellofresh/janus/pkg/opentracing"
	"github.com/hellofresh/janus/pkg/store"
	stats "github.com/hellofresh/stats-go"
	opentracing "github.com/opentracing/opentracing-go"
)

// initializes the basic configuration for the log wrapper
func initLogger(levelStr string) {
	level, err := log.ParseLevel(strings.ToLower(levelStr))
	if err != nil {
		log.WithError(err).Error("Error getting log level")
	}

	log.SetLevel(level)
	log.SetFormatter(&logrus_logstash.LogstashFormatter{
		Type:            "Janus",
		TimestampFormat: time.RFC3339Nano,
	})
}

// initializes distributed tracing
func initTracing(config config.Tracing) {
	log.Debug("initializing Open Tracing")
	tracer, err := tracerfactory.Build(config)
	if err != nil {
		log.WithError(err).Panic("Could not build a tracer for open tracing")
	}

	opentracing.InitGlobalTracer(tracer)
}

// initializes the storage and managers
func initStorage(config config.Storage) store.Store {
	storage, err := store.Build(config.DSN)
	if nil != err {
		log.Panic(err)
	}

	return storage
}

func initStatsdClient(config config.Stats) stats.Client {
	sectionsTestsMap, err := stats.ParseSectionsTestsMap(config.IDs)
	if err != nil {
		log.WithError(err).WithField("config", config.IDs).
			Error("Failed to parse stats second level IDs from env")
		sectionsTestsMap = map[stats.PathSection]stats.SectionTestDefinition{}
	}
	log.WithField("config", config.IDs).
		WithField("map", sectionsTestsMap.String()).
		Debug("Setting stats second level IDs")

	client, err := stats.NewClient(config.DSN, config.Prefix)
	if err != nil {
		log.WithError(err).Panic("Error initializing statsd client")
	}

	client.SetHTTPMetricCallback(stats.NewHasIDAtSecondLevelCallback(sectionsTestsMap))
	return client
}
