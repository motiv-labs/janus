package cors

import (
	"github.com/hellofresh/janus/pkg/plugin"
	"github.com/hellofresh/janus/pkg/proxy"
	"github.com/rs/cors"
)

type Config struct {
	AllowedOrigins []string `json:"domains"`
	AllowedMethods []string `json:"methods"`
	AllowedHeaders []string `json:"request_headers"`
	ExposedHeaders []string `json:"exposed_headers"`
}

func init() {
	plugin.RegisterPlugin("cors", plugin.Plugin{
		Action: setupCors,
	})
}

func setupCors(route *proxy.Route, p plugin.Params) error {
	var config Config

	err := plugin.Decode(p.Config, &config)
	if err != nil {
		return err
	}

	mw := cors.New(cors.Options{
		AllowedOrigins:   config.AllowedOrigins,
		AllowedMethods:   config.AllowedMethods,
		AllowedHeaders:   config.AllowedHeaders,
		ExposedHeaders:   config.ExposedHeaders,
		AllowCredentials: true,
	})

	route.AddInbound(mw.Handler)
	return nil
}
