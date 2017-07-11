package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_LoadEnv(t *testing.T) {
	os.Setenv("PORT", "8001")

	globalConfig, err := LoadEnv()
	require.NoError(t, err)

	assert.Equal(t, 8001, globalConfig.Port)
	assert.Equal(t, 8081, globalConfig.Web.Port)
}
