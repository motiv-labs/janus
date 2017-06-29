package oauth

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFileSystemRepository(t *testing.T) {
	wd, err := os.Getwd()
	assert.NoError(t, err)
	assert.Contains(t, wd, "github.com/hellofresh/janus")

	// .../github.com/hellofresh/janus/pkg/api/../../examples/auth
	exampleAPIsPath := filepath.Join(wd, "..", "..", "examples", "auth")
	info, err := os.Stat(exampleAPIsPath)
	assert.NoError(t, err)
	assert.True(t, info.IsDir())

	fsRepo, err := NewFileSystemRepository(exampleAPIsPath)
	assert.NoError(t, err)
	assert.NotNil(t, fsRepo)
}
