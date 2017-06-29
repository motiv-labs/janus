package compression

import (
	"github.com/hellofresh/janus/pkg/plugin"
	"github.com/hellofresh/janus/pkg/proxy"
	"github.com/pressly/chi/middleware"
)

func init() {
	plugin.RegisterPlugin("compression", plugin.Plugin{
		Action: setupCompression,
	})
}

func setupCompression(route *proxy.Route, p plugin.Params) error {
	route.AddInbound(middleware.DefaultCompress)
	return nil
}
