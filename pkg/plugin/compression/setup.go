package compression

import (
	"github.com/go-chi/chi/middleware"
	"github.com/hellofresh/janus/pkg/plugin"
	"github.com/hellofresh/janus/pkg/proxy"
)

func init() {
	plugin.RegisterPlugin("compression", plugin.Plugin{
		Action: setupCompression,
	})
}

func setupCompression(def *proxy.RouterDefinition, rawConfig plugin.Config) error {
	def.AddMiddleware(middleware.Compress(5))
	return nil
}
