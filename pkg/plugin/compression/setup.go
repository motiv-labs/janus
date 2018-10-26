package compression

import (
	"github.com/go-chi/chi/middleware"
	"github.com/hellofresh/janus/pkg/plugin"
	"github.com/hellofresh/janus/pkg/proxy"
)

func init() {
	plugin.RegisterPlugin("compression", plugin.Plugin{
		Action:   setupCompression,
		Validate: validateConfig,
	})
}

func setupCompression(def *proxy.RouterDefinition, rawConfig plugin.Config) error {
	def.AddMiddleware(middleware.DefaultCompress)
	return nil
}

func validateConfig(rawConfig plugin.Config) (bool, error) {
	return true, nil // This plugin does not have any configuration
}
