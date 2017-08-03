package responsetransformer

import (
	"github.com/hellofresh/janus/pkg/plugin"
	"github.com/hellofresh/janus/pkg/proxy"
)

func init() {
	plugin.RegisterPlugin("response_transformer", plugin.Plugin{
		Action: setupResponseTransformer,
	})
}

func setupResponseTransformer(route *proxy.Route, p plugin.Params) error {
	var config Config
	err := plugin.Decode(p.Config, &config)
	if err != nil {
		return err
	}

	route.AddOutbound(NewResponseTransformer(config))
	return nil
}
