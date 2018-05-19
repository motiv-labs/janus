package requesttransformer

import (
	"testing"

	"github.com/hellofresh/janus/pkg/api"
	"github.com/hellofresh/janus/pkg/plugin"
	"github.com/stretchr/testify/assert"
)

func TestRequestTransformerConfig(t *testing.T) {
	var config Config
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

	err := plugin.Decode(rawConfig, &config)
	assert.NoError(t, err)

	assert.IsType(t, map[string]string{}, config.Add.Headers)
	assert.Contains(t, config.Add.Headers, "NAME")

	assert.IsType(t, map[string]string{}, config.Add.QueryString)
	assert.Contains(t, config.Add.QueryString, "name")
}

func TestRequestTransformerPlugin(t *testing.T) {
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

	def := api.NewDefinition()
	err := setupRequestTransformer(def, rawConfig)
	assert.NoError(t, err)

	assert.Len(t, def.Proxy.Middleware(), 1)
}
