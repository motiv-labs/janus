package main

import (
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/bshuster-repo/logrus-logstash-hook"
	"github.com/hellofresh/janus/pkg/config"
	tracerfactory "github.com/hellofresh/janus/pkg/opentracing"
	"github.com/hellofresh/janus/pkg/store"
	opentracing "github.com/opentracing/opentracing-go"
)

var (
	err          error
	globalConfig *config.Specification
	storage      store.Store
)

// initializes the global configuration
func init() {
	conf, err := config.Load()
	globalConfig, err = config.Load()
	log.Debug(conf)
	if nil != err {
		log.WithError(err).Panic("Could not parse the environment configurations")
	}
}

// initializes the basic configuration for the log wrapper
func init() {
	level, err := log.ParseLevel(strings.ToLower(globalConfig.LogLevel))
	if err != nil {
		log.WithError(err).Error("Error getting log level")
	}

	log.SetLevel(level)
	log.SetFormatter(&logrus_logstash.LogstashFormatter{
		Type:            "Janus",
		TimestampFormat: time.RFC3339Nano,
	})
}

// initializes the
func init() {
	log.Debug("initializing Open Tracing")
	tracer, err := tracerfactory.Build(globalConfig.Tracing)
	if err != nil {
		log.WithError(err).Panic("Could not build a tracer for open tracing")
	}

	opentracing.InitGlobalTracer(tracer)
}

// initializes the storage and managers
func init() {
	var err error
	storage, err = store.Build(globalConfig.Storage.DSN)
	if nil != err {
		log.Panic(err)
	}
}
