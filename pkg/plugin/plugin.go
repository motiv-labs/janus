package plugin

import (
	"github.com/hellofresh/janus/pkg/api"
	"github.com/hellofresh/janus/pkg/router"
)

// Plugin defines basic methods for plugins
type Plugin interface {
	GetName() string
	GetMiddlewares(config api.Config, referenceSpec *api.Spec) ([]router.Constructor, error)
}

// Loader holds all availables plugins
type Loader struct {
	plugins map[string]Plugin
}

// NewLoader creates a new instance of Loader
func NewLoader() *Loader {
	return &Loader{
		plugins: make(map[string]Plugin),
	}
}

// Add a new plugin to the loader
func (l *Loader) Add(plugins ...Plugin) {
	for _, p := range plugins {
		l.plugins[p.GetName()] = p
	}
}

// Get a plugin by name
func (l *Loader) Get(name string) Plugin {
	return l.plugins[name]
}
