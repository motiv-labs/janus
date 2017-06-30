package plugin

import (
	"github.com/go-chi/chi/middleware"
	"github.com/hellofresh/janus/pkg/api"
	"github.com/hellofresh/janus/pkg/router"
)

// Compression represents the compression plugin
type Compression struct{}

// NewCompression creates a new instance of Compression
func NewCompression() *Compression {
	return &Compression{}
}

// GetName retrieves the plugin's name
func (h *Compression) GetName() string {
	return "compression"
}

// GetMiddlewares retrieves the plugin's middlewares
func (h *Compression) GetMiddlewares(rawConfig map[string]interface{}, referenceSpec *api.Spec) ([]router.Constructor, error) {
	return []router.Constructor{
		middleware.DefaultCompress,
	}, nil
}
