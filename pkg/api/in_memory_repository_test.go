package api

import (
	"testing"

	"github.com/hellofresh/janus/pkg/proxy"
	"github.com/stretchr/testify/assert"
)

func TestAPIInMemoryRepository(t *testing.T) {
	t.Parallel()

	tests := []struct {
		scenario string
		function func(*testing.T, Repository)
	}{
		{
			scenario: "api in memory exists",
			function: testExists,
		},
		{
			scenario: "api in memory not exists",
			function: testNotExists,
		},
		{
			scenario: "api in memory add missing name",
			function: testAddMissingName,
		},
		{
			scenario: "api in memory remove existent definition",
			function: testRemoveExistentDefinition,
		},
		{
			scenario: "api in memory remove inexistent definition",
			function: testRemoveInexistentDefinition,
		},
		{
			scenario: "api in memory find all",
			function: testFindAll,
		},
		{
			scenario: "api in memory find by listen path",
			function: testFindByListenPath,
		},
		{
			scenario: "api in memory not find by listen path",
			function: testNotFindByListenPath,
		},
		{
			scenario: "api in memory find by name",
			function: testFindByName,
		},
		{
			scenario: "api in memory not find by name",
			function: testNotFindByName,
		},
		{
			scenario: "api in memory not find by health checks",
			function: testFindValidHealthChecks,
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			repo := newInMemoryRepo()
			test.function(t, repo)
		})
	}
}

func testExists(t *testing.T, repo Repository) {
	definition := &Definition{
		Name: "test3",
		Proxy: &proxy.Definition{
			ListenPath: "/test3",
			Upstreams: &proxy.Upstreams{
				Targets: []*proxy.Target{
					{
						Target: "http://test3.com",
					},
				},
			},
		},
	}
	repo.Add(definition)

	ok, err := repo.Exists(definition)
	assert.Error(t, err)
	assert.True(t, ok)

	ok, err = repo.Exists(&Definition{Name: "Not valid", Proxy: proxy.NewDefinition()})
	assert.NoError(t, err)
	assert.False(t, ok)
}

func testNotExists(t *testing.T, repo Repository) {
	ok, err := repo.Exists(&Definition{Name: "Not valid", Proxy: proxy.NewDefinition()})
	assert.NoError(t, err)
	assert.False(t, ok)
}

func testAddMissingName(t *testing.T, repo Repository) {
	err := repo.Add(NewDefinition())
	assert.Error(t, err)
}

func testRemoveExistentDefinition(t *testing.T, repo Repository) {
	err := repo.Remove("test1")
	assert.NoError(t, err)
}

func testRemoveInexistentDefinition(t *testing.T, repo Repository) {
	err := repo.Remove("test")
	assert.Error(t, err)
}

func testFindAll(t *testing.T, repo Repository) {
	results, err := repo.FindAll()
	assert.NoError(t, err)
	assert.Len(t, results, 3)
}

func testFindByListenPath(t *testing.T, repo Repository) {
	definition, err := repo.FindByListenPath("/test1")
	assert.NoError(t, err)
	assert.NotNil(t, definition)
}

func testNotFindByListenPath(t *testing.T, repo Repository) {
	definition, err := repo.FindByListenPath("/invalid")
	assert.Error(t, err)
	assert.Nil(t, definition)
}

func testFindByName(t *testing.T, repo Repository) {
	definition, err := repo.FindByName("test1")
	assert.NoError(t, err)
	assert.NotNil(t, definition)
}

func testNotFindByName(t *testing.T, repo Repository) {
	definition, err := repo.FindByName("invalid")
	assert.Error(t, err)
	assert.Nil(t, definition)
}

func testFindValidHealthChecks(t *testing.T, repo Repository) {
	definitions, err := repo.FindValidAPIHealthChecks()
	assert.NoError(t, err)
	assert.Len(t, definitions, 1)
}

func newInMemoryRepo() *InMemoryRepository {
	repo := NewInMemoryRepository()

	repo.Add(&Definition{
		Name:   "test1",
		Active: true,
		Proxy: &proxy.Definition{
			ListenPath: "/test1",
			Upstreams: &proxy.Upstreams{
				Targets: []*proxy.Target{
					{
						Target: "http://test1.com",
					},
				},
			},
		},
		HealthCheck: HealthCheck{
			URL:     "http://test1.com/status.com",
			Timeout: 5,
		},
	})

	repo.Add(&Definition{
		Name:   "test2",
		Active: true,
		Proxy: &proxy.Definition{
			ListenPath: "/test2",
			Upstreams: &proxy.Upstreams{
				Targets: []*proxy.Target{
					{
						Target: "http://test2.com",
					},
				},
			},
		},
	})

	repo.Add(&Definition{
		Name: "test3",
		Proxy: &proxy.Definition{
			ListenPath: "/test3",
			Upstreams: &proxy.Upstreams{
				Targets: []*proxy.Target{
					{
						Target: "http://test3.com",
					},
				},
			},
		},
		HealthCheck: HealthCheck{
			URL:     "http://test3.com/status",
			Timeout: 5,
		},
	})

	return repo
}
