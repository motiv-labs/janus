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
	os.Setenv("BASIC_USERS", "admin:admin,test:test")
	os.Setenv("GITHUB_ORGANIZATIONS", "hellofresh,tests")
	os.Setenv("GITHUB_TEAMS", "hellofresh:tests,tests:devs")
	os.Setenv("JANUS_ADMIN_TEAM", "janus-owners")

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
	assert.Equal(t, map[string]string{"admin": "admin", "test": "test"}, globalConfig.Web.Credentials.Basic.Users)
	assert.Equal(t, map[string]string{"hellofresh": "tests", "tests": "devs"}, globalConfig.Web.Credentials.Github.Teams)
	assert.Equal(t, []string{"hellofresh", "tests"}, globalConfig.Web.Credentials.Github.Organizations)
	assert.Equal(t, "janus-owners", globalConfig.Web.Credentials.JanusAdminTeam)
	assert.True(t, globalConfig.Web.Credentials.Github.IsConfigured())
}
