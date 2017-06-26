package plugin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCORSConfig(t *testing.T) {
	var config corsConfig
	rawConfig := map[string]interface{}{
		"domains":         []string{"*"},
		"methods":         []string{"GET"},
		"request_headers": []string{"Content-Type", "Authorization"},
		"exposed_headers": []string{"Test"},
	}

	err := decode(rawConfig, &config)
	assert.NoError(t, err)

	assert.IsType(t, []string{}, config.AllowedOrigins)
	assert.Equal(t, []string{"*"}, config.AllowedOrigins)

	assert.IsType(t, []string{}, config.AllowedMethods)
	assert.Equal(t, []string{"GET"}, config.AllowedMethods)

	assert.IsType(t, []string{}, config.AllowedHeaders)
	assert.Equal(t, []string{"Content-Type", "Authorization"}, config.AllowedHeaders)

	assert.IsType(t, []string{}, config.ExposedHeaders)
	assert.Equal(t, []string{"Test"}, config.ExposedHeaders)
}

func TestInvalidCORSConfig(t *testing.T) {
	var config corsConfig
	rawConfig := map[string]interface{}{
		"domains": "*",
	}

	err := decode(rawConfig, &config)
	assert.Error(t, err)
}
