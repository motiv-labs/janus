package store_test

import (
	"testing"

	"github.com/hellofresh/janus/pkg/store"
	"github.com/stretchr/testify/assert"
)

func TestBuildInMemoryStore(t *testing.T) {
	storage, err := store.Build("memory://localhost")

	assert.Nil(t, err)
	assert.IsType(t, &store.InMemoryStore{}, storage)
	assert.Implements(t, (*store.Store)(nil), storage)
}

func TestWrongInMemoryDSN(t *testing.T) {
	_, err := store.Build("wrong://localhost")

	assert.NotNil(t, err)
	assert.IsType(t, store.ErrUnknownStorage, err)
}

func TestBuildNotParsableURL(t *testing.T) {
	_, err := store.Build("broken")

	assert.Error(t, err)
}

func TestBuildRedisPool(t *testing.T) {
	storage, err := store.Build("redis://localhost")

	assert.Nil(t, err)
	assert.IsType(t, &store.RedisStore{}, storage)
	assert.Implements(t, (*store.Store)(nil), storage)
}
