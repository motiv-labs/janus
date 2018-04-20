package compression

import (
	"github.com/go-chi/chi/middleware"
	"github.com/hellofresh/janus/pkg/api"
	"github.com/hellofresh/janus/pkg/plugin"
	"github.com/hellofresh/janus/pkg/proxy"
)

func init() {
	plugin.RegisterPlugin("compression", plugin.Plugin{
		Action: setupCompression,
	})
}

func setupCompression(def *api.Definition, route *proxy.Route, rawConfig plugin.Config) error {
	route.AddInbound(middleware.DefaultCompress)
	return nil
}
