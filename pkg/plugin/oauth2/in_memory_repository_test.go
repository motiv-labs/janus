package oauth2

import (
	"testing"

	"github.com/hellofresh/janus/pkg/proxy"
	"github.com/stretchr/testify/assert"
)

func TestOauthServerInMemoryRepository(t *testing.T) {
	t.Parallel()

	tests := []struct {
		scenario string
		function func(*testing.T, Repository)
	}{
		{
			scenario: "remove existent oauth server",
			function: testRemoveExistentOAuthServer,
		},
		{
			scenario: "remove inexistent oauth server",
			function: testRemoveInexistentOAuthServer,
		},
		{
			scenario: "find all oauth servers",
			function: testFindAllOAuthServers,
		},
		{
			scenario: "find by name",
			function: testFindByName,
		},
		{
			scenario: "not find by name",
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

func testRemoveExistentOAuthServer(t *testing.T, repo Repository) {
	err := repo.Remove("test1")
	assert.NoError(t, err)
}

func testRemoveInexistentOAuthServer(t *testing.T, repo Repository) {
	err := repo.Remove("test")
	assert.Error(t, err)
}

func testFindAllOAuthServers(t *testing.T, repo Repository) {
	results, err := repo.FindAll()
	assert.NoError(t, err)
	assert.Len(t, results, 2)
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

func newInMemoryRepo() *InMemoryRepository {
	repo := NewInMemoryRepository()
	repo.Add(&OAuth{
		Name: "test1",
		Endpoints: Endpoints{
			Token: &proxy.Definition{
				ListenPath: "/token",
				Upstreams: &proxy.Upstreams{
					Balancing: "roundrobin",
					Targets: []*proxy.Target{
						&proxy.Target{Target: "http://test.com/token"},
					},
				},
			},
		},
	})

	repo.Add(&OAuth{
		Name: "test2",
		Endpoints: Endpoints{
			Token: &proxy.Definition{
				ListenPath: "/token",
				Upstreams: &proxy.Upstreams{
					Balancing: "roundrobin",
					Targets: []*proxy.Target{
						&proxy.Target{Target: "http://test2.com/token"},
					},
				},
			},
		},
	})

	return repo
}
