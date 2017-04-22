package store

import (
	"encoding/json"

	log "github.com/Sirupsen/logrus"
	"github.com/garyburd/redigo/redis"
	"github.com/hellofresh/janus/pkg/notifier"
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
		Prefix: "janus",
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

// Exists checks if a key exists in the store
func (s *RedisStore) Exists(key string) (bool, error) {
	conn := s.getConnection()
	defer conn.Close()

	return s.exists(conn, key)
}

// Get retreives a value from the store
func (s *RedisStore) Get(key string) (string, error) {
	conn := s.getConnection()
	defer conn.Close()

	return s.get(conn, key)
}

// Remove a value from the store
func (s *RedisStore) Remove(key string) error {
	conn := s.getConnection()
	defer conn.Close()

	return s.remove(conn, key)
}

// Set a value in the store
func (s *RedisStore) Set(key string, value string, expire int64) error {
	conn := s.getConnection()
	defer conn.Close()

	return s.set(conn, key, value, expire)
}

// Publish publishes to a topic in redis
func (s *RedisStore) Publish(topic string, data []byte) error {
	c := s.getConnection()
	_, err := c.Do("PUBLISH", topic, data)
	return err
}

// Subscribe subscribes to a topic in redis
func (s *RedisStore) Subscribe(channel string, callback func(notifier.Notification)) error {
	// Get a connection from a pool
	c := s.getConnection()
	defer c.Close()

	psc := redis.PubSubConn{Conn: c}
	if err := psc.Subscribe(channel); err != nil {
		return err
	}

	for {
		switch v := psc.Receive().(type) {
		case redis.Message:
			notif := notifier.Notification{}
			if err := json.Unmarshal(v.Data, &notif); err != nil {
				log.Error("Unmarshalling message body failed, malformed: ", err)
				return err
			}
			callback(notif)
		case error:
			log.WithError(v).Debug("An error ocurred when getting the message")
			return v
		}
	}
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
	}

	args = append(args, key)
	args = append(args, expire)
	args = append(args, value)
	return "SETEX", args
}
