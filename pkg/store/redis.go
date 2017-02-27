package store

import (
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

func (s *RedisStore) getConnection() redis.Conn {
	return s.Pool.Get()
}

// ping checks if redis is alive.
func (s *RedisStore) ping() (bool, error) {
	conn := s.getConnection()
	defer conn.Close()

	data, err := conn.Do("PING")
	if err != nil || data == nil {
		return false, err
	}

	return data == "PONG", nil
}

func (s *RedisStore) Exists(key string) (bool, error) {
	conn := s.getConnection()
	defer conn.Close()

	return s.exists(conn, key)
}

func (s *RedisStore) Get(key string) (string, error) {
	conn := s.getConnection()
	defer conn.Close()

	return s.get(conn, key)
}

func (s *RedisStore) Remove(key string) error {
	conn := s.getConnection()
	defer conn.Close()

	return s.remove(conn, key)
}

func (s *RedisStore) Set(key string, value string, expire int64) error {
	conn := s.getConnection()
	defer conn.Close()

	return s.set(conn, key, value, expire)
}

func (s *RedisStore) ToLimiterStore(prefix string) (limiter.Store, error) {
	// Alternatively, you can pass options to the store with the "WithOptions"
	// function. For example, for Redis store:
	return limiter.NewRedisStoreWithOptions(s.Pool, limiter.StoreOptions{
		Prefix:   prefix,
		MaxRetry: limiter.DefaultMaxRetry,
	})
}

func (s *RedisStore) exists(conn redis.Conn, key string) (bool, error) {
	exists, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (s *RedisStore) remove(conn redis.Conn, key string) error {
	_, err := conn.Do("DEL", key)
	return err
}

func (s *RedisStore) get(conn redis.Conn, key string) (string, error) {
	return redis.String(conn.Do("GET", key))
}

func (s *RedisStore) set(conn redis.Conn, key string, value string, expire int64) error {
	command, args := getSetCommandAndArgs(key, value, expire)
	if _, err := conn.Do(command, args...); err != nil {
		return err
	}

	return nil
}

func getSetCommandAndArgs(key string, value string, expire int64) (string, []interface{}) {
	var args []interface{}
	if expire == 0 {
		args = append(args, key)
		args = append(args, value)
		return "SET", args
	} else {
		args = append(args, key)
		args = append(args, expire)
		args = append(args, value)
		return "SETEX", args
	}
}
