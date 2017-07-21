package bodylmt

import (
	"github.com/hellofresh/janus/pkg/plugin"
	"github.com/hellofresh/janus/pkg/proxy"
)

// Config represents the Body Limit configuration
type Config struct {
	Limit string `json:"limit"`
}

func init() {
	plugin.RegisterPlugin("body_limit", plugin.Plugin{
		Action: setupBodyLimit,
	})
}

func setupBodyLimit(route *proxy.Route, p plugin.Params) error {
	var config Config
	err := plugin.Decode(p.Config, &config)
	if err != nil {
		return err
	}

	route.AddInbound(NewBodyLimitMiddleware(config.Limit))
	return nil
}
