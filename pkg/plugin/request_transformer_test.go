package plugin

import (
	"testing"

	"github.com/hellofresh/janus/pkg/middleware"
	"github.com/stretchr/testify/assert"
)

func TestRequestTransformerConfig(t *testing.T) {
	var config middleware.RequestTransformerConfig
	rawConfig := map[string]interface{}{
		"add": map[string]interface{}{
			"headers": map[string]string{
				"NAME": "TEST",
			},
			"querystring": map[string]string{
				"name": "test",
			},
		},
	}

	err := decode(rawConfig, &config)
	assert.NoError(t, err)

	assert.IsType(t, map[string]string{}, config.Add.Headers)
	assert.Contains(t, config.Add.Headers, "NAME")

	assert.IsType(t, map[string]string{}, config.Add.QueryString)
	assert.Contains(t, config.Add.QueryString, "name")
}
