package store

import (
	"net/url"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/garyburd/redigo/redis"
)

const (
	// InMemory storage
	InMemory = "memory"
	// Redis storage
	Redis = "redis"
	// None Nullable storage
	None = "none"
)

// Store is the common interface for datastores.
type Store interface {
	Exists(key string) (bool, error)
	Get(key string) (string, error)
	Remove(key string) error
	Set(key string, value string, expire int64) error
}

// Options are options for store.
type Options struct {
	// Prefix is the prefix to use for the key.
	Prefix string

	// MaxRetry is the maximum number of retry under race conditions.
	MaxRetry int

	// CleanUpInterval is the interval for cleanup.
	CleanUpInterval time.Duration
}

// Build creates a new storage based on the provided DSN
func Build(dsn string) (Store, error) {
	url, err := url.Parse(dsn)
	if nil != err {
		return nil, err
	}
	log.WithField("type", url.Scheme).Debug("Initializing storage")

	switch url.Scheme {
	case InMemory:
		return NewInMemoryStore(), nil
	case Redis:
		// Create a Redis pool.
		pool := &redis.Pool{
			MaxIdle:     3,
			IdleTimeout: 240 * time.Second,
			Dial:        func() (redis.Conn, error) { return redis.DialURL(dsn) },
		}

		log.WithField("dsn", dsn).Debug("Trying to connect to redis pool")
		return NewRedisStore(pool, url.Query().Get("prefix"))
	}

	return nil, ErrUnknownStorage
}
