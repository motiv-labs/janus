package api

import (
	"testing"

	"github.com/hellofresh/janus/pkg/proxy"
	"github.com/stretchr/testify/assert"
)

func newInMemoryRepo() *InMemoryRepository {
	repo := NewInMemoryRepository()

	repo.Add(&Definition{
		Name: "test1",
		Proxy: &proxy.Definition{
			ListenPath:  "/test1",
			UpstreamURL: "http://test1",
		},
		HealthCheck: HealthCheck{
			URL:     "http://test1.com/status",
			Timeout: 5,
		},
	})

	repo.Add(&Definition{
		Name: "test2",
		Proxy: &proxy.Definition{
			ListenPath:  "/test2",
			UpstreamURL: "http://test2",
		},
	})

	return repo
}

func TestExists(t *testing.T) {
	definition := &Definition{
		Name: "test3",
		Proxy: &proxy.Definition{
			ListenPath:  "/test3",
			UpstreamURL: "http://test3",
		},
	}
	repo := newInMemoryRepo()
	repo.Add(definition)

	ok, err := repo.Exists(definition)
	assert.Error(t, err)
	assert.True(t, ok)

	ok, err = repo.Exists(&Definition{Name: "Not valid", Proxy: proxy.NewDefinition()})
	assert.NoError(t, err)
	assert.False(t, ok)
}

func TestNotExists(t *testing.T) {
	repo := newInMemoryRepo()

	ok, err := repo.Exists(&Definition{Name: "Not valid", Proxy: proxy.NewDefinition()})
	assert.NoError(t, err)
	assert.False(t, ok)
}

func TestAddMissingName(t *testing.T) {
	repo := newInMemoryRepo()

	err := repo.Add(NewDefinition())
	assert.Error(t, err)
}

func TestRemoveExistentDefinition(t *testing.T) {
	repo := newInMemoryRepo()

	err := repo.Remove("test1")
	assert.NoError(t, err)
}

func TestRemoveInexistentDefinition(t *testing.T) {
	repo := newInMemoryRepo()

	err := repo.Remove("test")
	assert.Error(t, err)
}

func TestFindAll(t *testing.T) {
	repo := newInMemoryRepo()

	results, err := repo.FindAll()
	assert.NoError(t, err)
	assert.Len(t, results, 2)
}

func TestFindByListenPath(t *testing.T) {
	repo := newInMemoryRepo()

	definition, err := repo.FindByListenPath("/test1")
	assert.NoError(t, err)
	assert.NotNil(t, definition)
}

func TestNotFindByListenPath(t *testing.T) {
	repo := newInMemoryRepo()

	definition, err := repo.FindByListenPath("/invalid")
	assert.Error(t, err)
	assert.Nil(t, definition)
}

func TestFindByName(t *testing.T) {
	repo := newInMemoryRepo()

	definition, err := repo.FindByName("test1")
	assert.NoError(t, err)
	assert.NotNil(t, definition)
}

func TestNotFindByName(t *testing.T) {
	repo := newInMemoryRepo()

	definition, err := repo.FindByName("invalid")
	assert.Error(t, err)
	assert.Nil(t, definition)
}

func TestFindValidHealthChecks(t *testing.T) {
	repo := newInMemoryRepo()

	definitions, err := repo.FindValidAPIHealthChecks()
	assert.NoError(t, err)
	assert.Len(t, definitions, 1)
}
