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

	// go-redis Client instance
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
	data, err := s.Client.Ping().Result()
	return data == "PONG", err
}

// Exists checks if a key exists in the store
func (s *RedisStore) Exists(key string) (bool, error) {
	data, err := s.Client.Exists(key).Result()
	return data == 1, err
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
func (s *RedisStore) Publish(topic string, data []byte) error {
	return s.Client.Publish(topic, data).Err()
}

// Subscribe subscribes to a topic in redis
func (s *RedisStore) Subscribe(channel string, callback func(notifier.Notification)) error {
	pubsub := s.Client.Subscribe(channel)
	defer pubsub.Close()

	for {
		v, err := pubsub.ReceiveMessage()
		if err != nil {
			log.WithError(err).Debug("An error occurred when getting the message")
			return err
		}
		notification := notifier.Notification{}
		if marshallErr := json.Unmarshal([]byte(v.Payload), &notification); marshallErr != nil {
			log.WithError(marshallErr).Error("Unmarshalling message body failed, malformed")
			return marshallErr
		}
		callback(notification)
	}
}
