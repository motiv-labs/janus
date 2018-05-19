package requesttransformer

import (
	"github.com/hellofresh/janus/pkg/api"
	"github.com/hellofresh/janus/pkg/plugin"
)

func init() {
	plugin.RegisterPlugin("request_transformer", plugin.Plugin{
		Action: setupRequestTransformer,
	})
}

func setupRequestTransformer(def *api.Definition, rawConfig plugin.Config) error {
	var config Config
	err := plugin.Decode(rawConfig, &config)
	if err != nil {
		return err
	}

	def.Proxy.AddMiddleware(NewRequestTransformer(config))
	return nil
}
