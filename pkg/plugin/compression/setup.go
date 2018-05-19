package compression

import (
	"github.com/go-chi/chi/middleware"
	"github.com/hellofresh/janus/pkg/api"
	"github.com/hellofresh/janus/pkg/plugin"
)

func init() {
	plugin.RegisterPlugin("compression", plugin.Plugin{
		Action: setupCompression,
	})
}

func setupCompression(def *api.Definition, rawConfig plugin.Config) error {
	def.Proxy.AddMiddleware(middleware.DefaultCompress)
	return nil
}
