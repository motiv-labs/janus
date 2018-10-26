package responsetransformer

import (
	"github.com/hellofresh/janus/pkg/plugin"
	"github.com/hellofresh/janus/pkg/proxy"
)

func init() {
	plugin.RegisterPlugin("response_transformer", plugin.Plugin{
		Action:   setupResponseTransformer,
		Validate: validateConfig,
	})
}

func setupResponseTransformer(def *proxy.RouterDefinition, rawConfig plugin.Config) error {
	var config Config
	err := plugin.Decode(rawConfig, &config)
	if err != nil {
		return err
	}

	def.AddMiddleware(NewResponseTransformer(config))
	return nil
}

func validateConfig(rawConfig plugin.Config) (bool, error) {
	return true, nil // This plugin does not have any configuration
}
