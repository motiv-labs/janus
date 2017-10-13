package store

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis"
	"github.com/hellofresh/janus/pkg/notifier"
	log "github.com/sirupsen/logrus"
)

const defaultPrefix = "janus"

// RedisStore is the redis store.
type RedisStore struct {
	// The prefix to use for the key.
	Prefix string

	// github.com/go-redis/redis Client instance.
	Client *redis.Client

	// The maximum number of retry under race conditions.
	MaxRetry int
}

// NewRedisStore returns an instance of redis store.
func NewRedisStore(client *redis.Client, prefix string) (Store, error) {
	if prefix == "" {
		prefix = defaultPrefix
	}

	return NewRedisStoreWithOptions(client, Options{
		Prefix: prefix,
	})
}

// NewRedisStoreWithOptions returns an instance of redis store with custom options.
func NewRedisStoreWithOptions(client *redis.Client, options Options) (Store, error) {
	store := &RedisStore{
		Client: client,
		Prefix: options.Prefix,
	}

	return store, nil
}

// ping checks if redis is alive.
func (s *RedisStore) ping() (bool, error) {
	val, err := s.Client.Ping().Result()
	if err != nil {
		return false, err
	}

	return val == "PONG", nil
}

// Exists checks if a key exists in the store
func (s *RedisStore) Exists(key string) (bool, error) {
	val, err := s.Client.Exists(key).Result()
	if err != nil {
		return false, err
	}

	return val > 0, nil
}

// Get retreives a value from the store
func (s *RedisStore) Get(key string) (string, error) {
	return s.Client.Get(key).Result()
}

// Remove a value from the store
func (s *RedisStore) Remove(key string) error {
	return s.Client.Del(key).Err()
}

// Set a value in the store
func (s *RedisStore) Set(key string, value string, expire int64) error {
	return s.Client.Set(key, value, time.Duration(expire)*time.Second).Err()
}

// Publish publishes to a topic in redis
func (s *RedisStore) Publish(channel string, data []byte) error {
	return s.Client.Publish(channel, string(data)).Err()
}

// Subscribe subscribes to a topic in redis
func (s *RedisStore) Subscribe(channel string, callback func(notifier.Notification)) error {
	pubsub := s.Client.Subscribe(channel)
	defer pubsub.Close()

	for {
		msg, err := pubsub.ReceiveMessage()
		if err != nil {
			return err
		}

		notification := notifier.Notification{}
		if err := json.Unmarshal([]byte(msg.Payload), &notification); err != nil {
			log.WithError(err).Error("Unmarshalling message body failed, malformed")
			return err
		}
		callback(notification)
	}
}
