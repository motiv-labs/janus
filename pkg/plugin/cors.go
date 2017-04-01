package plugin

import (
	"github.com/hellofresh/janus/pkg/api"
	"github.com/hellofresh/janus/pkg/router"
	"github.com/rs/cors"
)

// CORS represents the cors plugin
type CORS struct{}

// NewCORS creates a new instance of CORS
func NewCORS() *CORS {
	return &CORS{}
}

// GetName retrieves the plugin's name
func (h *CORS) GetName() string {
	return "cors"
}

// GetMiddlewares retrieves the plugin's middlewares
func (h *CORS) GetMiddlewares(config api.Config, referenceSpec *api.Spec) ([]router.Constructor, error) {
	middleware := cors.New(cors.Options{
		AllowedOrigins:   convertToSlice(config["domains"]),
		AllowedMethods:   convertToSlice(config["methods"]),
		AllowedHeaders:   convertToSlice(config["request_headers"]),
		ExposedHeaders:   convertToSlice(config["exposed_headers"]),
		AllowCredentials: true,
	})

	return []router.Constructor{
		middleware.Handler,
	}, nil
}

func convertToSlice(config interface{}) []string {
	var values []string
	aInterface := config.([]interface{})
	for _, v := range aInterface {
		values = append(values, v.(string))
	}

	return values
}
