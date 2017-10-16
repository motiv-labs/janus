package store

import (
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type RedisTestSuite struct {
	suite.Suite
	redisServer *miniredis.Miniredis
	redisStore  Store
}

func (suite *RedisTestSuite) SetupTest() {
	server, err := miniredis.Run()
	suite.Nil(err)
	suite.redisServer = server

	client := redis.NewClient(&redis.Options{
		Addr:        suite.redisServer.Addr(),
		DB:          0,
		PoolSize:    3,
		IdleTimeout: 240 * time.Second,
	})
	store, err := NewRedisStore(client, "testPrefix")
	suite.Nil(err)
	suite.redisStore = store
}

func (suite *RedisTestSuite) TearDownTest() {
	suite.redisStore = nil
	suite.redisServer.Close()
	suite.redisServer = nil
}

func (suite *RedisTestSuite) TestSet() {
	assert.Nil(suite.T(), suite.redisStore.Set(testKey, testValue, 0))
	suite.redisServer.CheckGet(suite.T(), testKey, testValue)
}

func (suite *RedisTestSuite) TestSetEx() {
	assert.Nil(suite.T(), suite.redisStore.Set(testKey, testValue, 20))
	suite.redisServer.CheckGet(suite.T(), testKey, testValue)
	// Still exists after 10 seconds
	suite.redisServer.FastForward(10 * time.Second)
	assert.True(suite.T(), suite.redisServer.Exists(testKey))
	suite.redisServer.CheckGet(suite.T(), testKey, testValue)
	// Doesn't exists after 21 seconds (total)
	suite.redisServer.FastForward(11 * time.Second)
	assert.False(suite.T(), suite.redisServer.Exists(testKey))
}

func (suite *RedisTestSuite) TestGet() {
	suite.redisServer.Set(testKey, testValue)
	value, err := suite.redisStore.Get(testKey)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), value, testValue)

	_, err = suite.redisStore.Get("DoesntExists")
	assert.NotNil(suite.T(), err)
}

func (suite *RedisTestSuite) TestExists() {
	// Key doesn't exist yet
	exists, err := suite.redisStore.Exists(testKey)
	assert.Nil(suite.T(), err)
	assert.False(suite.T(), exists)

	// Set the key and check existance
	suite.redisServer.Set(testKey, testValue)
	exists, err = suite.redisStore.Exists(testKey)
	assert.Nil(suite.T(), err)
	assert.True(suite.T(), exists)

	// Try closing the server and see we receive an error
	suite.redisServer.Close()
	exists, err = suite.redisStore.Exists(testKey)
	assert.NotNil(suite.T(), err)
	assert.False(suite.T(), exists)
}

func (suite *RedisTestSuite) TestRemove() {
	suite.redisServer.Set(testKey, testValue)
	// Check it exists
	value, err := suite.redisStore.Get(testKey)
	assert.Equal(suite.T(), value, testValue)
	assert.Nil(suite.T(), err)

	// Able to remove
	assert.Nil(suite.T(), suite.redisStore.Remove(testKey))

	// Doesn't exist anymore
	value, err = suite.redisStore.Get(testKey)
	assert.Empty(suite.T(), value)
	assert.NotNil(suite.T(), err)
}

func TestRedisTestSuite(t *testing.T) {
	suite.Run(t, new(RedisTestSuite))
}
