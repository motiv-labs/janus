package cors

import (
	"github.com/asaskevich/govalidator"
	"github.com/hellofresh/janus/pkg/plugin"
	"github.com/hellofresh/janus/pkg/proxy"
	"github.com/rs/cors"
)

// Config represents the CORS configuration
type Config struct {
	AllowedOrigins     []string `json:"domains"`
	AllowedMethods     []string `json:"methods"`
	AllowedHeaders     []string `json:"request_headers"`
	ExposedHeaders     []string `json:"exposed_headers"`
	OptionsPassthrough bool     `json:"options_passthrough"`
}

func init() {
	plugin.RegisterPlugin("cors", plugin.Plugin{
		Action:   setupCors,
		Validate: validateConfig,
	})
}

func setupCors(def *proxy.RouterDefinition, rawConfig plugin.Config) error {
	var config Config

	err := plugin.Decode(rawConfig, &config)
	if err != nil {
		return err
	}

	mw := cors.New(cors.Options{
		AllowedOrigins:     config.AllowedOrigins,
		AllowedMethods:     config.AllowedMethods,
		AllowedHeaders:     config.AllowedHeaders,
		ExposedHeaders:     config.ExposedHeaders,
		OptionsPassthrough: config.OptionsPassthrough,
		AllowCredentials:   true,
	})

	def.AddMiddleware(mw.Handler)
	return nil
}

func validateConfig(rawConfig plugin.Config) (bool, error) {
	var config Config
	err := plugin.Decode(rawConfig, &config)
	if err != nil {
		return false, err
	}

	return govalidator.ValidateStruct(config)
}
