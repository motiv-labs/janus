package requesttransformer

import (
	"github.com/hellofresh/janus/pkg/plugin"
	"github.com/hellofresh/janus/pkg/proxy"
)

func init() {
	plugin.RegisterPlugin("request_transformer", plugin.Plugin{
		Action: setupRequestTransformer,
	})
}

func setupRequestTransformer(route *proxy.Route, p plugin.Params) error {
	var config Config
	err := plugin.Decode(p.Config, &config)
	if err != nil {
		return err
	}

	route.AddInbound(NewRequestTransformer(config))
	return nil
}
