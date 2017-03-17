package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ulule/limiter"
)

const (
	testKey   = "key"
	testValue = "value"
)

func TestInMemoryStore_NewInMemoryStore(t *testing.T) {
	instance := NewInMemoryStore()

	assert.IsType(t, &InMemoryStore{}, instance)
	assert.Implements(t, (*Store)(nil), instance)
}

func TestInMemoryStore_Set(t *testing.T) {
	instance := NewInMemoryStore()

	err := instance.Set(testKey, testValue, 0)
	assert.Nil(t, err)
}

func TestInMemoryStore_Exists(t *testing.T) {
	instance := NewInMemoryStore()

	instance.Set(testKey, testValue, 0)

	val, err := instance.Exists("foo")
	assert.Nil(t, err)
	assert.False(t, val)

	val, err = instance.Exists(testKey)
	assert.Nil(t, err)
	assert.True(t, val)
}

func TestInMemoryStore_Get(t *testing.T) {
	instance := NewInMemoryStore()

	instance.Set(testKey, testValue, 0)

	val, err := instance.Get(testKey)
	assert.Nil(t, err)
	assert.Equal(t, testValue, val)

	val, err = instance.Get("foo")
	assert.Nil(t, err)
	assert.Empty(t, val)
}

func TestInMemoryStore_Remove(t *testing.T) {
	instance := NewInMemoryStore()

	instance.Set(testKey, testValue, 0)

	val, err := instance.Get(testKey)
	assert.Nil(t, err)
	assert.Equal(t, testValue, val)

	instance.Remove(testKey)

	val, err = instance.Get(testKey)
	assert.Nil(t, err)
	assert.Empty(t, val)
}

func TestInMemoryStore_ToLimiterStore(t *testing.T) {
	instance := NewInMemoryStore()

	limiterStore, err := instance.ToLimiterStore("prefix")
	assert.Nil(t, err)
	assert.IsType(t, &limiter.MemoryStore{}, limiterStore)
}
