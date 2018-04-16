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
		function func(*testing.T, *InMemoryRepository)
	}{
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
			scenario: "api in memory find by name",
			function: testFindByName,
		},
		{
			scenario: "api in memory not find by name",
			function: testNotFindByName,
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			repo := newInMemoryRepo()
			test.function(t, repo)
		})
	}
}

func testAddMissingName(t *testing.T, repo *InMemoryRepository) {
	err := repo.add(NewDefinition())
	assert.Error(t, err)
}

func testRemoveExistentDefinition(t *testing.T, repo *InMemoryRepository) {
	err := repo.remove("test1")
	assert.NoError(t, err)
}

func testRemoveInexistentDefinition(t *testing.T, repo *InMemoryRepository) {
	err := repo.remove("test")
	assert.Error(t, err)
}

func testFindAll(t *testing.T, repo *InMemoryRepository) {
	results, err := repo.FindAll()
	assert.NoError(t, err)
	assert.Len(t, results, 3)
}

func testFindByName(t *testing.T, repo *InMemoryRepository) {
	definition, err := repo.findByName("test1")
	assert.NoError(t, err)
	assert.NotNil(t, definition)
}

func testNotFindByName(t *testing.T, repo *InMemoryRepository) {
	definition, err := repo.findByName("invalid")
	assert.Error(t, err)
	assert.Nil(t, definition)
}

func newInMemoryRepo() *InMemoryRepository {
	repo := NewInMemoryRepository()

	repo.add(&Definition{
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

	repo.add(&Definition{
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

	repo.add(&Definition{
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
