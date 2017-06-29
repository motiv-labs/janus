package oauth

import (
	"testing"

	"net/url"

	"github.com/hellofresh/janus/pkg/proxy"
	"github.com/stretchr/testify/assert"
)

func newInMemoryRepo() *InMemoryRepository {
	repo := NewInMemoryRepository()

	repo.Add(&OAuth{
		Name: "test1",
		Endpoints: Endpoints{
			Token: &proxy.Definition{
				ListenPath:  "/token",
				UpstreamURL: "http://test.com/token",
			},
		},
	})

	repo.Add(&OAuth{
		Name: "test2",
		Endpoints: Endpoints{
			Token: &proxy.Definition{
				ListenPath:  "/token",
				UpstreamURL: "http://test2.com/token",
			},
		},
	})

	return repo
}

func TestRemoveExistentOAuthServer(t *testing.T) {
	repo := newInMemoryRepo()

	err := repo.Remove("test1")
	assert.NoError(t, err)
}

func TestRemoveInexistentOAuthServer(t *testing.T) {
	repo := newInMemoryRepo()

	err := repo.Remove("test")
	assert.Error(t, err)
}

func TestFindAllOAuthServers(t *testing.T) {
	repo := newInMemoryRepo()

	results, err := repo.FindAll()
	assert.NoError(t, err)
	assert.Len(t, results, 2)
}

func TestFindByTokenURL(t *testing.T) {
	repo := newInMemoryRepo()

	tokenURL, err := url.Parse("http://test.com/token")
	assert.NoError(t, err)

	server, err := repo.FindByTokenURL(*tokenURL)
	assert.NoError(t, err)
	assert.NotNil(t, server)
}

func TestNotFindByTokenURL(t *testing.T) {
	repo := newInMemoryRepo()

	tokenURL, err := url.Parse("http://test.com/wrong")
	assert.NoError(t, err)

	server, err := repo.FindByTokenURL(*tokenURL)
	assert.Error(t, err)
	assert.Nil(t, server)
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
