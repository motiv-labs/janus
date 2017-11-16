package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testKey   = "key"
	testValue = "value"
)

func TestInMemoryStore(t *testing.T) {
	t.Parallel()

	tests := []struct {
		scenario string
		function func(*testing.T, *InMemoryStore)
	}{
		{
			scenario: "new in memory store",
			function: testNewInMemoryStore,
		},
		{
			scenario: "set in memory store",
			function: testInMemoryStoreSet,
		},
		{
			scenario: "check exists from in memory store",
			function: testInMemoryStoreExists,
		},
		{
			scenario: "get from in memory store",
			function: testInMemoryStoreGet,
		},
		{
			scenario: "remove from in memory store",
			function: testInMemoryStoreRemove,
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			instance := NewInMemoryStore()
			test.function(t, instance)
		})
	}
}

func testNewInMemoryStore(t *testing.T, instance *InMemoryStore) {
	assert.IsType(t, &InMemoryStore{}, instance)
	assert.Implements(t, (*Store)(nil), instance)
}

func testInMemoryStoreSet(t *testing.T, instance *InMemoryStore) {
	err := instance.Set(testKey, testValue, 0)
	assert.Nil(t, err)
}

func testInMemoryStoreExists(t *testing.T, instance *InMemoryStore) {
	instance.Set(testKey, testValue, 0)

	val, err := instance.Exists("foo")
	assert.Nil(t, err)
	assert.False(t, val)

	val, err = instance.Exists(testKey)
	assert.Nil(t, err)
	assert.True(t, val)
}

func testInMemoryStoreGet(t *testing.T, instance *InMemoryStore) {
	instance.Set(testKey, testValue, 0)

	val, err := instance.Get(testKey)
	assert.Nil(t, err)
	assert.Equal(t, testValue, val)

	val, err = instance.Get("foo")
	assert.Nil(t, err)
	assert.Empty(t, val)
}

func testInMemoryStoreRemove(t *testing.T, instance *InMemoryStore) {
	instance.Set(testKey, testValue, 0)

	val, err := instance.Get(testKey)
	assert.Nil(t, err)
	assert.Equal(t, testValue, val)

	instance.Remove(testKey)

	val, err = instance.Get(testKey)
	assert.Nil(t, err)
	assert.Empty(t, val)
}
