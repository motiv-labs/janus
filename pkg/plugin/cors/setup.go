package cors

import (
	"github.com/hellofresh/janus/pkg/plugin"
	"github.com/hellofresh/janus/pkg/proxy"
	"github.com/rs/cors"
)

// Config represents the CORS configuration
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

func setupCors(def *proxy.RouterDefinition, rawConfig plugin.Config) error {
	var config Config

	err := plugin.Decode(rawConfig, &config)
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

	def.AddMiddleware(mw.Handler)
	return nil
}
