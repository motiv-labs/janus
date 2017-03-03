package store_test

import (
	"testing"

	"github.com/hellofresh/janus/pkg/store"
	"github.com/stretchr/testify/assert"
	"github.com/ulule/limiter"
)

var (
	storage = store.NewInMemoryStore()
)

func TestNewInMemoryStore(t *testing.T) {
	assert.IsType(t, &store.InMemoryStore{}, storage)
	assert.Implements(t, (*store.Store)(nil), storage)
}

func TestSetValueToStore(t *testing.T) {
	err := storage.Set("key", "test", 0)
	assert.Nil(t, err)
}

func TestGetValueFromStore(t *testing.T) {
	value, err := storage.Get("key")
	assert.Nil(t, err)

	assert.Equal(t, "test", value)
}

func TestKeyExistis(t *testing.T) {
	exists, err := storage.Exists("key")
	assert.Nil(t, err)

	assert.True(t, exists)
}

func TestRemove(t *testing.T) {
	err := storage.Set("keyToRemove", "test", 0)
	assert.Nil(t, err)

	exists, err := storage.Exists("keyToRemove")
	assert.Nil(t, err)
	assert.True(t, exists)

	err = storage.Remove("keyToRemove")
	assert.Nil(t, err)

	exists, err = storage.Exists("keyToRemove")
	assert.Nil(t, err)
	assert.False(t, exists)
}

func TestConvertToLmiterStore(t *testing.T) {
	limiterStore, err := storage.ToLimiterStore("prefix")
	assert.Nil(t, err)
	assert.IsType(t, &limiter.MemoryStore{}, limiterStore)
}
