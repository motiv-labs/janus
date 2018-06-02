package bodylmt

import (
	"github.com/hellofresh/janus/pkg/api"
	"github.com/hellofresh/janus/pkg/plugin"
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

func setupBodyLimit(def *api.Definition, rawConfig plugin.Config) error {
	var config Config
	err := plugin.Decode(rawConfig, &config)
	if err != nil {
		return err
	}

	def.Proxy.AddMiddleware(NewBodyLimitMiddleware(config.Limit))
	return nil
}
