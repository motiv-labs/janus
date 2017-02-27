package store

import (
	"errors"
	"testing"

	"github.com/garyburd/redigo/redis"
	"github.com/stretchr/testify/assert"
	"github.com/ulule/limiter"
)

type mockConnDoAssert struct {
	t        *testing.T
	command  string
	args     []interface{}
	doResult interface{}
	doError  error
}

func (conn mockConnDoAssert) Do(command string, args ...interface{}) (interface{}, error) {
	assert.Equal(conn.t, conn.command, command)
	assert.Len(conn.t, args, len(conn.args))
	for idx, val := range args {
		assert.Equal(conn.t, conn.args[idx], val)
	}
	return conn.doResult, conn.doError
}
func (conn mockConnDoAssert) Send(string, ...interface{}) error { return nil }
func (conn mockConnDoAssert) Err() error                        { return nil }
func (conn mockConnDoAssert) Close() error                      { return nil }
func (conn mockConnDoAssert) Flush() error                      { return nil }
func (conn mockConnDoAssert) Receive() (interface{}, error)     { return nil, nil }

type mockRedisStore struct {
	conn       redis.Conn
	redisStore *RedisStore
}

func newMockRedisStore(conn redis.Conn) *mockRedisStore {
	return &mockRedisStore{conn: conn, redisStore: &RedisStore{}}
}

func (s *mockRedisStore) getConnection() redis.Conn {
	return s.conn
}

func (s *mockRedisStore) Exists(key string) (bool, error) {
	conn := s.getConnection()
	defer conn.Close()

	return s.redisStore.exists(conn, key)
}

func (s *mockRedisStore) Get(key string) (string, error) {
	conn := s.getConnection()
	defer conn.Close()

	return s.redisStore.get(conn, key)
}

func (s *mockRedisStore) Remove(key string) error {
	conn := s.getConnection()
	defer conn.Close()

	return s.redisStore.remove(conn, key)
}

func (s *mockRedisStore) Set(key string, value string, expire int64) error {
	conn := s.getConnection()
	defer conn.Close()

	return s.redisStore.set(conn, key, value, expire)
}

func (s *mockRedisStore) ToLimiterStore(prefix string) (limiter.Store, error) {
	// TODO: implement tests for limiter
	return nil, nil
}

func TestRedisStore_getSetCommandAndArgs(t *testing.T) {
	command, args := getSetCommandAndArgs(testKey, testValue, 0)
	assert.Equal(t, "SET", command)
	assert.Len(t, args, 2)
	assert.Equal(t, testKey, args[0])
	assert.Equal(t, testValue, args[1])

	seconds := 3600
	command, args = getSetCommandAndArgs(testKey, testValue, int64(seconds))
	assert.Equal(t, "SETEX", command)
	assert.Len(t, args, 3)
	assert.Equal(t, testKey, args[0])
	assert.Equal(t, int64(seconds), args[1])
	assert.Equal(t, testValue, args[2])
}

func TestRedisStore_Set(t *testing.T) {
	var argsSet []interface{}
	argsSet = append(argsSet, testKey)
	argsSet = append(argsSet, testValue)

	connectionSet := mockConnDoAssert{t: t, command: "SET", args: argsSet, doResult: nil, doError: nil}
	instanceSet := newMockRedisStore(connectionSet)
	assert.Nil(t, instanceSet.Set(testKey, testValue, 0))

	seconds := int64(3600)
	var argsSetEx []interface{}
	argsSetEx = append(argsSetEx, testKey)
	argsSetEx = append(argsSetEx, seconds)
	argsSetEx = append(argsSetEx, testValue)

	connectionSetEx := mockConnDoAssert{t: t, command: "SETEX", args: argsSetEx, doResult: nil, doError: nil}
	instanceSetEx := newMockRedisStore(connectionSetEx)
	assert.Nil(t, instanceSetEx.Set(testKey, testValue, seconds))
}

func TestInMemoryStore_Set_week(t *testing.T) {
	seconds := int64(2629743)
	var argsSetEx []interface{}
	argsSetEx = append(argsSetEx, testKey)
	argsSetEx = append(argsSetEx, seconds)
	argsSetEx = append(argsSetEx, testValue)

	connectionSetEx := mockConnDoAssert{t: t, command: "SETEX", args: argsSetEx, doResult: nil, doError: nil}
	instanceSetEx := newMockRedisStore(connectionSetEx)
	assert.Nil(t, instanceSetEx.Set(testKey, testValue, seconds))
}

func TestRedisStore_Get(t *testing.T) {
	var argsSet []interface{}
	argsSet = append(argsSet, testKey)
	connection := mockConnDoAssert{t: t, command: "GET", args: argsSet, doResult: testValue, doError: nil}
	instance := newMockRedisStore(connection)

	val, err := instance.Get(testKey)
	assert.Nil(t, err)
	assert.Equal(t, testValue, val)
}

func TestRedisStore_Exists(t *testing.T) {
	var argsSet []interface{}
	argsSet = append(argsSet, testKey)

	connectionTrue := mockConnDoAssert{t: t, command: "EXISTS", args: argsSet, doResult: int64(1), doError: nil}
	instanceTrue := newMockRedisStore(connectionTrue)

	val, err := instanceTrue.Exists(testKey)
	assert.Nil(t, err)
	assert.True(t, val)

	connectionFalse := mockConnDoAssert{t: t, command: "EXISTS", args: argsSet, doResult: int64(0), doError: nil}
	instanceFalse := newMockRedisStore(connectionFalse)
	val, err = instanceFalse.Exists(testKey)
	assert.Nil(t, err)
	assert.False(t, val)

	connectionError := mockConnDoAssert{t: t, command: "EXISTS", args: argsSet, doResult: int64(0), doError: errors.New("")}
	instanceError := newMockRedisStore(connectionError)
	val, err = instanceError.Exists(testKey)
	assert.NotNil(t, err)
	assert.False(t, val)
}

func TestRedisStore_Remove(t *testing.T) {
	var argsRemove []interface{}
	argsRemove = append(argsRemove, testKey)
	connection := mockConnDoAssert{t: t, command: "DEL", args: argsRemove, doError: nil}
	instance := newMockRedisStore(connection)

	err := instance.Remove(testKey)
	assert.Nil(t, err)
}
