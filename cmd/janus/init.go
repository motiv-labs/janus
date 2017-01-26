package main

import (
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/bshuster-repo/logrus-logstash-hook"
	"github.com/garyburd/redigo/redis"
	"github.com/hellofresh/janus/pkg/config"
	"github.com/hellofresh/janus/pkg/middleware"
	"github.com/hellofresh/janus/pkg/store"
	statsd "gopkg.in/alexcesaro/statsd.v2"
)

var (
	err          error
	globalConfig *config.Specification
	accessor     *middleware.DatabaseAccessor
	storage      store.Store
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
		Type:            globalConfig.Application.Name,
		TimestampFormat: time.RFC3339Nano,
	})
}

// initializes a DB connection
func init() {
	accessor, err = middleware.InitDB(globalConfig.Database.DSN)
	if err != nil {
		log.Fatalf("Couldn't connect to the mongodb database: %s", err.Error())
	}
}

// initializes a storage
func init() {
	// Create a Redis pool.
	pool := &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", globalConfig.StorageDSN) },
	}

	dsn := globalConfig.StorageDSN
	log.Debugf("Trying to connect to redis pool: %s", dsn)
	storage, err = store.NewRedisStore(pool)
	if err != nil {
		log.Fatalf("Couldn't connect to the redis pool: %s", err.Error())
	}
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
