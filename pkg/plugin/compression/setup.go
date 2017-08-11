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

func setupCompression(route *proxy.Route, rawConfig plugin.Config) error {
	route.AddInbound(middleware.DefaultCompress)
	return nil
}
