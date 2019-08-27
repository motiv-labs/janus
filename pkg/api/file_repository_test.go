package api

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/hellofresh/janus/pkg/proxy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newRepo(t *testing.T) *FileSystemRepository {
	wd, err := os.Getwd()
	assert.NoError(t, err)

	// ./../../assets/apis
	exampleAPIsPath := filepath.Join(wd, "..", "..", "assets", "apis")
	info, err := os.Stat(exampleAPIsPath)
	require.NoError(t, err)
	require.True(t, info.IsDir())

	fsRepo, err := NewFileSystemRepository(exampleAPIsPath)
	assert.NoError(t, err)

	return fsRepo
}

func TestNewFileSystemRepository(t *testing.T) {
	fsRepo := newRepo(t)

	allDefinitions, err := fsRepo.FindAll()
	assert.NoError(t, err)
	assert.Equal(t, 3, len(allDefinitions))

	defToAdd := &Definition{
		Name: "foo-bar",
		Proxy: &proxy.Definition{
			ListenPath: "/foo/bar/*",
			Upstreams: &proxy.Upstreams{
				Balancing: "roundrobin",
				Targets: []*proxy.Target{
					{Target: "http://example.com/foo/bar/"},
				},
			},
		},
	}
	err = fsRepo.add(defToAdd)
	require.NoError(t, err)

	def, err := fsRepo.findByName(defToAdd.Name)
	require.NoError(t, err)
	assert.Equal(t, defToAdd.Name, def.Name)
	assert.Equal(t, defToAdd.Proxy.ListenPath, def.Proxy.ListenPath)
}

func assertFindByName(t *testing.T, fsRepo *FileSystemRepository) {
	def, err := fsRepo.findByName("example")
	assert.NoError(t, err)
	assert.Equal(t, "example", def.Name)
	assert.Equal(t, "/example/*", def.Proxy.ListenPath)

	_, err = fsRepo.findByName("foo")
	assert.Equal(t, ErrAPIDefinitionNotFound, err)
}

func TestFileSystemRepository_Add(t *testing.T) {
	fsRepo := newRepo(t)

	invalidName := &Definition{Name: ""}
	err := fsRepo.add(invalidName)
	assert.Error(t, err)
}
