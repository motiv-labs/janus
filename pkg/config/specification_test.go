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
	os.Setenv("BASIC_USERS", "admin:admin, test:test")
	os.Setenv("GITHUB_ORGANIZATIONS", "hellofresh, tests")
	os.Setenv("GITHUB_TEAMS", "hellofresh:tests, hellofresh:devs")

	globalConfig, err := LoadEnv()
	require.NoError(t, err)

	assert.Equal(t, 8001, globalConfig.Port)
	assert.Equal(t, 8081, globalConfig.Web.Port)
	assert.Equal(t, "HS256", globalConfig.Web.Credentials.Algorithm)
	assert.Equal(t, uint(0), globalConfig.Stats.AutoDiscoverThreshold)
	assert.Equal(t, []string{"api", "foo", "bar"}, globalConfig.Stats.AutoDiscoverWhiteList)
	assert.Equal(t, "error-log", globalConfig.Stats.ErrorsSection)
	assert.False(t, globalConfig.TLS.IsHTTPS())
	assert.False(t, globalConfig.Tracing.IsGoogleCloudEnabled())
	assert.False(t, globalConfig.Tracing.IsAppdashEnabled())
	assert.IsType(t, map[string]string{}, globalConfig.Web.Credentials.Basic.Users)
	assert.Len(t, globalConfig.Web.Credentials.Basic.Users, 2)
	assert.IsType(t, map[string]string{}, globalConfig.Web.Credentials.Github.Teams)
	assert.Len(t, globalConfig.Web.Credentials.Github.Teams, 2)
	assert.IsType(t, []string{}, globalConfig.Web.Credentials.Github.Organizations)
	assert.Len(t, globalConfig.Web.Credentials.Github.Organizations, 2)
	assert.True(t, globalConfig.Web.Credentials.Github.IsConfigured())
}
