package oauth2

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFileSystemRepository(t *testing.T) {
	wd, err := os.Getwd()
	assert.NoError(t, err)

	// ./../../assets/auth
	exampleAPIsPath := filepath.Join(wd, "..", "..", "..", "assets", "stubs", "auth")
	info, err := os.Stat(exampleAPIsPath)
	require.NoError(t, err)
	require.True(t, info.IsDir())

	fsRepo, err := NewFileSystemRepository(exampleAPIsPath)
	require.NoError(t, err)
	assert.NotNil(t, fsRepo)
}
