package config

import (
	"os"
	"testing"
	"time"

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
	os.Setenv("TOKEN_TIMEOUT", "2h")

	globalConfig, err := LoadEnv()
	require.NoError(t, err)

	assert.Equal(t, 8001, globalConfig.Port)
	assert.Equal(t, 8081, globalConfig.Web.Port)
	assert.Equal(t, "HS256", globalConfig.Web.Credentials.Algorithm)
	assert.Equal(t, uint(0), globalConfig.Stats.AutoDiscoverThreshold)
	assert.Equal(t, []string{"api", "foo", "bar"}, globalConfig.Stats.AutoDiscoverWhiteList)
	assert.Equal(t, "error-log", globalConfig.Stats.ErrorsSection)
	assert.False(t, globalConfig.TLS.IsHTTPS())
	assert.Equal(t, map[string]string{"admin": "admin", "test": "test"}, globalConfig.Web.Credentials.Basic.Users)
	assert.Equal(t, map[string]string{"hellofresh": "tests", "tests": "devs"}, globalConfig.Web.Credentials.Github.Teams)
	assert.Equal(t, []string{"hellofresh", "tests"}, globalConfig.Web.Credentials.Github.Organizations)
	assert.Equal(t, "janus-owners", globalConfig.Web.Credentials.JanusAdminTeam)
	assert.Equal(t, 2*time.Hour, globalConfig.Web.Credentials.Timeout)
	assert.True(t, globalConfig.Web.Credentials.Github.IsConfigured())

}

func TestDefaults(t *testing.T) {
	os.Clearenv()
	globalConfig, err := LoadEnv()
	require.NoError(t, err)

	assert.Equal(t, 8080, globalConfig.Port)
	assert.Equal(t, 8081, globalConfig.Web.Port)
	assert.Equal(t, 8433, globalConfig.TLS.Port)
	assert.Equal(t, 8444, globalConfig.Web.TLS.Port)
	assert.True(t, globalConfig.Web.TLS.Redirect)
	assert.Equal(t, 20*time.Millisecond, globalConfig.BackendFlushInterval)
	assert.Equal(t, 180*time.Second, globalConfig.RespondingTimeouts.IdleTimeout)
	assert.True(t, globalConfig.TLS.Redirect)
	assert.True(t, globalConfig.RequestID)
	assert.Equal(t, 10*time.Second, globalConfig.Cluster.UpdateFrequency)
	assert.Equal(t, "file:///etc/janus", globalConfig.Database.DSN)

	assert.Equal(t, "HS256", globalConfig.Web.Credentials.Algorithm)
	assert.Equal(t, map[string]string{"admin": "admin"}, globalConfig.Web.Credentials.Basic.Users)

	assert.Equal(t, "log://", globalConfig.Stats.DSN)
	assert.Equal(t, "error-log", globalConfig.Stats.ErrorsSection)

	assert.Equal(t, "janus", globalConfig.Tracing.ServiceName)
	assert.Equal(t, "always", globalConfig.Tracing.SamplingStrategy)
	assert.Equal(t, 0.15, globalConfig.Tracing.SamplingParam)
}
