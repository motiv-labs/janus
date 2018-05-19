package responsetransformer

import (
	"github.com/hellofresh/janus/pkg/api"
	"github.com/hellofresh/janus/pkg/plugin"
)

func init() {
	plugin.RegisterPlugin("response_transformer", plugin.Plugin{
		Action: setupResponseTransformer,
	})
}

func setupResponseTransformer(def *api.Definition, rawConfig plugin.Config) error {
	var config Config
	err := plugin.Decode(rawConfig, &config)
	if err != nil {
		return err
	}

	def.Proxy.AddMiddleware(NewResponseTransformer(config))
	return nil
}
