package plugin

import (
	"encoding/json"
	"testing"

	"github.com/asaskevich/govalidator"
	"github.com/hellofresh/janus/pkg/api"
	"github.com/stretchr/testify/assert"
)

type TestPluginConfig struct {
	Number int    `json:"number"`
	Text   string `json:"text"`
}

func init() {
	RegisterPlugin("test_plugin", Plugin{
		Validate: validateConfig,
	})
	RegisterPlugin("test_plugin_without_validation_function", Plugin{})
}

func validateConfig(rawConfig Config) (bool, error) {
	var config TestPluginConfig
	err := Decode(rawConfig, &config)
	if err != nil {
		return false, err
	}

	return govalidator.ValidateStruct(config)
}

func TestSuccessfulPluginValidation(t *testing.T) {
	definition := api.NewDefinition()
	json.Unmarshal([]byte(validDefinition), definition)

	for _, v := range definition.Plugins {
		result, err := ValidateConfig(v.Name, v.Config)
		assert.NoError(t, err)
		assert.True(t, result)
	}
}

func TestFailedPluginValidation(t *testing.T) {
	definition := api.NewDefinition()
	json.Unmarshal([]byte(invalidDefinition), definition)

	for _, v := range definition.Plugins {
		result, err := ValidateConfig(v.Name, v.Config)
		assert.Error(t, err)
		assert.False(t, result)
	}
}

func TestPluginValidationWithoutValidationFunction(t *testing.T) {
	definition := api.NewDefinition()
	json.Unmarshal([]byte(noValidationFuncDefinition), definition)

	for _, v := range definition.Plugins {
		result, err := ValidateConfig(v.Name, v.Config)
		assert.NoError(t, err)
		assert.True(t, result)
	}
}

const (
	validDefinition = `{
    "name" : "users",
    "active" : true,
    "proxy" : {
        "preserve_host" : false,
        "listen_path" : "/users/*",
        "upstreams" : {
            "balancing" : "weight",
            "targets" : [ 
                {
                    "target" : "http://localhost:8000/users",
                    "weight" : 0
                }, 
                {
                    "target" : "http://auth-service.live-k8s.hellofresh.io/users",
                    "weight" : 100
                }
            ]
        },
        "insecure_skip_verify" : false,
        "strip_path" : true,
        "append_path" : false,
        "enable_load_balancing" : false,
        "methods" : [ 
            "ALL"
        ],
        "hosts" : []
    },
    "plugins" : [
		{
			"name": "test_plugin",
			"enabled": true,
			"config": {
				"number": 420,
				"text": "Lorem ipsum dolor sit amet"
			}
		}
    ],
    "health_check" : {
        "url" : "",
        "timeout" : 0
    }
}`
	invalidDefinition = `{
    "name" : "users",
    "active" : true,
    "proxy" : {
        "preserve_host" : false,
        "listen_path" : "/users/*",
        "upstreams" : {
            "balancing" : "weight",
            "targets" : [ 
                {
                    "target" : "http://localhost:8000/users",
                    "weight" : 0
                }, 
                {
                    "target" : "http://auth-service.live-k8s.hellofresh.io/users",
                    "weight" : 100
                }
            ]
        },
        "insecure_skip_verify" : false,
        "strip_path" : true,
        "append_path" : false,
        "enable_load_balancing" : false,
        "methods" : [ 
            "ALL"
        ],
        "hosts" : []
    },
    "plugins" : [ 
		{
			"name": "test_plugin",
			"enabled": true,
			"config": {
				"number": "Not a number",
				"text": "Lorem ipsum dolor sit amet"
			}
		}
    ],
    "health_check" : {
        "url" : "",
        "timeout" : 0
    }
}`
	noValidationFuncDefinition = `{
    "name" : "users",
    "active" : true,
    "proxy" : {
        "preserve_host" : false,
        "listen_path" : "/users/*",
        "upstreams" : {
            "balancing" : "weight",
            "targets" : [ 
                {
                    "target" : "http://localhost:8000/users",
                    "weight" : 0
                }, 
                {
                    "target" : "http://auth-service.live-k8s.hellofresh.io/users",
                    "weight" : 100
                }
            ]
        },
        "insecure_skip_verify" : false,
        "strip_path" : true,
        "append_path" : false,
        "enable_load_balancing" : false,
        "methods" : [ 
            "ALL"
        ],
        "hosts" : []
    },
    "plugins" : [ 
		{
			"name": "test_plugin_without_validation_function",
			"enabled": true,
			"config": {
				"number": "Not a number",
				"text": "Lorem ipsum dolor sit amet"
			}
		}
    ],
    "health_check" : {
        "url" : "",
        "timeout" : 0
    }
}`
)
