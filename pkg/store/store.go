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

// Subscriber holds the basic methods to subscribe to a topic
type Subscriber interface {
	Subscribe(topic string) *Subscription
}

// Publisher holds the basic methods to publish a message
type Publisher interface {
	Publish(topic string, data []byte) error
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

// Message represents the message that comes
// form an update
type Message []byte

// Subscription holds a message channel
type Subscription struct {
	Message chan Message
}

// NewSubscription creates a new instance of Subscription
func NewSubscription() *Subscription {
	return &Subscription{
		Message: make(chan Message),
	}
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

		log.Debugf("Trying to connect to redis pool: %s", dsn)
		return NewRedisStore(pool)
	}

	return nil, ErrUnknownStorage
}
