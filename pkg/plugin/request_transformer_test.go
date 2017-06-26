package plugin

import (
	"testing"

	"github.com/hellofresh/janus/pkg/api"
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

func TestRequestTransformerPluginGetName(t *testing.T) {
	plugin := NewRequestTransformer()
	assert.Equal(t, "request_transformer", plugin.GetName())
}

func TestRequestTransformerPluginLocalPolicy(t *testing.T) {
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

	spec := &api.Spec{
		Definition: &api.Definition{
			Name: "API Name",
		},
	}

	plugin := NewRequestTransformer()
	middleware, err := plugin.GetMiddlewares(rawConfig, spec)

	assert.NoError(t, err)
	assert.Len(t, middleware, 1)
}
