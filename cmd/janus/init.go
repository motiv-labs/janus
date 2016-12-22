package main

import (
	"strings"

	statsd "gopkg.in/alexcesaro/statsd.v2"
	redis "gopkg.in/redis.v3"

	log "github.com/Sirupsen/logrus"
	"github.com/bshuster-repo/logrus-logstash-hook"
	"github.com/hellofresh/janus/config"
	"github.com/hellofresh/janus/middleware"
)

var (
	err          error
	globalConfig *config.Specification
	accessor     *middleware.DatabaseAccessor
	redisStorage *redis.Client
	statsdClient *statsd.Client
)

// initializes the global configuration
func init() {
	globalConfig, err = config.LoadEnv()
	if nil != err {
		log.Panic(err.Error())
	}
}

// initializes the basic configuration for the log wrapper
func init() {
	level, err := log.ParseLevel(strings.ToLower(globalConfig.LogLevel))
	if err != nil {
		log.Error("Error getting level", err)
	}

	log.SetLevel(level)
	log.SetFormatter(&logrus_logstash.LogstashFormatter{
		Type: globalConfig.Application.Name,
	})
}

// initializes a DB connection
func init() {
	accessor, err = middleware.InitDB(globalConfig.Database.DSN)
	if err != nil {
		log.Fatalf("Couldn't connect to the mongodb database: %s", err.Error())
	}
}

// initializes a Redis connection
func init() {
	dsn := globalConfig.StorageDSN
	log.Debugf("Trying to connect to redis instance: %s", dsn)
	redisStorage = redis.NewClient(&redis.Options{
		Addr: dsn,
	})
}

// initializes new StatsD client if it enabled
func init() {
	var options []statsd.Option
	statsdConfig := globalConfig.Statsd

	log.Debugf("Trying to connect to statsd instance: %s", statsdConfig.DSN)
	if len(statsdConfig.DSN) == 0 {
		log.Debug("Statsd DSN not provided, client will be muted")
		options = append(options, statsd.Mute(true))
	} else {
		options = append(options, statsd.Address(statsdConfig.DSN))
	}

	if len(statsdConfig.Prefix) > 0 {
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
