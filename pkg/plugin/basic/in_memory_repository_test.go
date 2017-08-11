package basic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func newInMemoryRepo() *InMemoryRepository {
	repo := NewInMemoryRepository()

	repo.Add(&User{
		Username: "test1",
		Password: "test1",
	})

	repo.Add(&User{
		Username: "test2",
		Password: "test2",
	})

	return repo
}

func TestAdd(t *testing.T) {
	repo := newInMemoryRepo()

	err := repo.Add(&User{
		Username: "test3",
		Password: "test3",
	})
	assert.NoError(t, err)
}

func TestRemoveExistentUser(t *testing.T) {
	repo := newInMemoryRepo()

	err := repo.Remove("test1")
	assert.NoError(t, err)
}

func TestRemoveNonexistentUser(t *testing.T) {
	repo := newInMemoryRepo()

	err := repo.Remove("invalid")
	assert.Error(t, err)
}

func TestFindAll(t *testing.T) {
	repo := newInMemoryRepo()

	results, err := repo.FindAll()
	assert.NoError(t, err)
	assert.Len(t, results, 2)
}

func TestFindByUsername(t *testing.T) {
	repo := newInMemoryRepo()

	result, err := repo.FindByUsername("test1")
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestNotFindByUsername(t *testing.T) {
	repo := newInMemoryRepo()

	result, err := repo.FindByUsername("invalid")
	assert.Error(t, err)
	assert.Nil(t, result)
}
