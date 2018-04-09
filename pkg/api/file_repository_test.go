package api

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/hellofresh/janus/pkg/proxy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newRepo(t *testing.T) *FileSystemRepository {
	wd, err := os.Getwd()
	assert.NoError(t, err)
	assert.Contains(t, wd, "github.com/hellofresh/janus")

	// .../github.com/hellofresh/janus/pkg/api/../../assets/apis
	exampleAPIsPath := filepath.Join(wd, "..", "..", "assets", "apis")
	info, err := os.Stat(exampleAPIsPath)
	assert.NoError(t, err)
	assert.True(t, info.IsDir())

	fsRepo, err := NewFileSystemRepository(exampleAPIsPath)
	assert.NoError(t, err)

	return fsRepo
}

func TestNewFileSystemRepository(t *testing.T) {
	fsRepo := newRepo(t)

	allDefinitions, err := fsRepo.FindAll()
	assert.NoError(t, err)
	assert.Equal(t, 3, len(allDefinitions))

	healthDefinitions, err := fsRepo.FindValidAPIHealthChecks()
	assert.NoError(t, err)
	assert.Equal(t, 3, len(healthDefinitions))

	assertFindByName(t, fsRepo)
	assertFindByFindByListenPath(t, fsRepo)
	assertExists(t, fsRepo)

	defToAdd := &Definition{
		Name: "foo-bar",
		Proxy: &proxy.Definition{
			ListenPath: "/foo/bar/*",
			Upstreams: &proxy.Upstreams{
				Balancing: "roundrobin",
				Targets: []*proxy.Target{
					&proxy.Target{Target: "http://example.com/foo/bar/"},
				},
			},
		},
	}
	err = fsRepo.Add(defToAdd)
	require.NoError(t, err)

	def, err := fsRepo.FindByName(defToAdd.Name)
	require.NoError(t, err)
	assert.Equal(t, defToAdd.Name, def.Name)
	assert.Equal(t, defToAdd.Proxy.ListenPath, def.Proxy.ListenPath)

	def, err = fsRepo.FindByListenPath(defToAdd.Proxy.ListenPath)
	require.NoError(t, err)
	assert.Equal(t, defToAdd.Name, def.Name)
	assert.Equal(t, defToAdd.Proxy.ListenPath, def.Proxy.ListenPath)

	exists, err := fsRepo.Exists(&Definition{Name: defToAdd.Name})
	assert.True(t, exists)
	assert.Equal(t, ErrAPINameExists, err)

	exists, err = fsRepo.Exists(&Definition{
		Name:  time.Now().Format(time.RFC3339Nano),
		Proxy: &proxy.Definition{ListenPath: defToAdd.Proxy.ListenPath},
	})
	assert.True(t, exists)
	assert.Equal(t, ErrAPIListenPathExists, err)

	err = fsRepo.Remove(defToAdd.Name)
	require.NoError(t, err)

	err = fsRepo.Remove(defToAdd.Name)
	assert.Equal(t, ErrAPIDefinitionNotFound, err)

	_, err = fsRepo.FindByName(defToAdd.Name)
	assert.Equal(t, ErrAPIDefinitionNotFound, err)

	_, err = fsRepo.FindByListenPath(defToAdd.Proxy.ListenPath)
	assert.Equal(t, ErrAPIDefinitionNotFound, err)

	exists, err = fsRepo.Exists(defToAdd)
	require.NoError(t, err)
	assert.False(t, exists)
}

func assertFindByName(t *testing.T, fsRepo *FileSystemRepository) {
	def, err := fsRepo.FindByName("example")
	assert.NoError(t, err)
	assert.Equal(t, "example", def.Name)
	assert.Equal(t, "/example/*", def.Proxy.ListenPath)

	_, err = fsRepo.FindByName("foo")
	assert.Equal(t, ErrAPIDefinitionNotFound, err)
}

func assertFindByFindByListenPath(t *testing.T, fsRepo *FileSystemRepository) {
	def, err := fsRepo.FindByListenPath("/example/*")
	assert.NoError(t, err)
	assert.Equal(t, "example", def.Name)
	assert.Equal(t, "/example/*", def.Proxy.ListenPath)

	_, err = fsRepo.FindByListenPath("/foo/*")
	assert.Equal(t, ErrAPIDefinitionNotFound, err)
}

func assertExists(t *testing.T, fsRepo *FileSystemRepository) {
	exists, err := fsRepo.Exists(&Definition{Name: "example"})
	assert.True(t, exists)
	assert.Equal(t, ErrAPINameExists, err)

	exists, err = fsRepo.Exists(&Definition{Name: "example1", Proxy: &proxy.Definition{ListenPath: "/example/*"}})
	assert.True(t, exists)
	assert.Equal(t, ErrAPIListenPathExists, err)

	exists, err = fsRepo.Exists(&Definition{Name: "example1", Proxy: &proxy.Definition{ListenPath: "/example1/*"}})
	assert.False(t, exists)
	assert.NoError(t, err)
}

func TestFileSystemRepository_Add(t *testing.T) {
	fsRepo := newRepo(t)

	invalidName := &Definition{Name: ""}
	err := fsRepo.Add(invalidName)
	assert.Error(t, err)
}
