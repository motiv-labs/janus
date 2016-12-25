package store

import (
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/ulule/limiter"
)

const (
	// DefaultPrefix is the default prefix to use for the key in the store.
	DefaultPrefix = "limiter"
)

// RedisStore is the redis store.
type RedisStore struct {
	// The prefix to use for the key.
	Prefix string

	// github.com/garyburd/redigo Pool instance.
	Pool *redis.Pool

	// The maximum number of retry under race conditions.
	MaxRetry int
}

// NewRedisStore returns an instance of redis store.
func NewRedisStore(pool *redis.Pool) (Store, error) {
	return NewRedisStoreWithOptions(pool, Options{
		Prefix: DefaultPrefix,
	})
}

// NewRedisStoreWithOptions returns an instance of redis store with custom options.
func NewRedisStoreWithOptions(pool *redis.Pool, options Options) (Store, error) {
	store := &RedisStore{
		Pool:   pool,
		Prefix: options.Prefix,
	}

	if _, err := store.ping(); err != nil {
		return nil, err
	}

	return store, nil
}

// ping checks if redis is alive.
func (s *RedisStore) ping() (bool, error) {
	conn := s.Pool.Get()
	defer conn.Close()

	data, err := conn.Do("PING")
	if err != nil || data == nil {
		return false, err
	}

	return data == "PONG", nil
}

func (s *RedisStore) Exists(key string) (bool, error) {
	conn := s.Pool.Get()
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (s *RedisStore) Get(key string) (string, error) {
	conn := s.Pool.Get()
	defer conn.Close()
	return redis.String(conn.Do("GET", key))
}

func (s *RedisStore) Set(key string, value string, expire time.Duration) error {
	conn := s.Pool.Get()
	defer conn.Close()

	_, err := conn.Do("SETEX", key, expire.Seconds(), value)
	if err != nil {
		return err
	}

	return nil
}

func (s *RedisStore) ToLimiterStore(prefix string) (limiter.Store, error) {
	// Alternatively, you can pass options to the store with the "WithOptions"
	// function. For example, for Redis store:
	return limiter.NewRedisStoreWithOptions(s.Pool, limiter.StoreOptions{
		Prefix:   prefix,
		MaxRetry: limiter.DefaultMaxRetry,
	})
}
