package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_LoadEnv(t *testing.T) {
	os.Setenv("PORT", "8001")
	os.Setenv("STATS_AUTO_DISCOVER_WHITE_LIST", "api,foo,bar")

	globalConfig, err := LoadEnv()
	require.NoError(t, err)

	assert.Equal(t, 8001, globalConfig.Port)
	assert.Equal(t, 8081, globalConfig.Web.Port)
	assert.Equal(t, uint(0), globalConfig.Stats.AutoDiscoverThreshold)
	assert.Equal(t, []string{"api", "foo", "bar"}, globalConfig.Stats.AutoDiscoverWhiteList)
	assert.Equal(t, "error-log", globalConfig.Stats.ErrorsSection)
	assert.False(t, globalConfig.TLS.IsHTTPS())
	assert.False(t, globalConfig.Tracing.IsGoogleCloudEnabled())
	assert.False(t, globalConfig.Tracing.IsAppdashEnabled())
}
