package oauth2

import (
	"testing"

	"github.com/hellofresh/janus/pkg/plugin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOAuth2Config(t *testing.T) {
	var config Config
	rawConfig := map[string]interface{}{
		"server_name": "test",
	}

	err := plugin.Decode(rawConfig, &config)
	require.NoError(t, err)
	assert.Equal(t, "test", config.ServerName)
}
