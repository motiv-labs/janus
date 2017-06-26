package plugin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOAuth2Config(t *testing.T) {
	var config oauth2Config
	rawConfig := map[string]interface{}{
		"server_name": "test",
	}

	err := decode(rawConfig, &config)
	assert.NoError(t, err)
	assert.Equal(t, "test", config.ServerName)
}
