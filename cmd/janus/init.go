package main

import (
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/bshuster-repo/logrus-logstash-hook"
	"github.com/hellofresh/janus/pkg/config"
	"github.com/hellofresh/janus/pkg/store"
	statsd "gopkg.in/alexcesaro/statsd.v2"
)

var (
	err          error
	globalConfig *config.Specification
	storage      store.Store
	statsdClient *statsd.Client
)

// initializes the global configuration
func init() {
	globalConfig, err = config.LoadEnv()
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
		Type:            globalConfig.Application.Name,
		TimestampFormat: time.RFC3339Nano,
	})
}

// initializes the storage and managers
func init() {
	var err error
	storage, err = store.Build(globalConfig.StorageDSN)
	if nil != err {
		log.Panic(err)
	}
}

// initializes new StatsD client if it enabled
func init() {
	var options []statsd.Option
	statsdConfig := globalConfig.Statsd

	log.Debugf("Trying to connect to statsd instance: %s", statsdConfig.DSN)
	if statsdConfig.IsEnabled() {
		log.Debug("Statsd DSN not provided, client will be muted")
		options = append(options, statsd.Mute(true))
	} else {
		options = append(options, statsd.Address(statsdConfig.DSN))
	}

	if statsdConfig.HasPrefix() {
		options = append(options, statsd.Prefix(statsdConfig.Prefix))
	}

	statsdClient, err = statsd.New(options...)
	if err != nil {
		log.WithError(err).
			WithFields(log.Fields{
				"dsn":    statsdConfig.DSN,
				"prefix": statsdConfig.Prefix,
			}).Warning("An error occurred while connecting to StatsD. Client will be muted.")
	}
}
