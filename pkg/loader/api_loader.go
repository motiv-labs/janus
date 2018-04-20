package loader

import (
	"github.com/hellofresh/janus/pkg/api"
	"github.com/hellofresh/janus/pkg/middleware"
	"github.com/hellofresh/janus/pkg/plugin"
	"github.com/hellofresh/janus/pkg/proxy"
	log "github.com/sirupsen/logrus"
)

// APILoader is responsible for loading all apis form a datastore and configure them in a register
type APILoader struct {
	register *proxy.Register
}

// NewAPILoader creates a new instance of the api manager
func NewAPILoader(register *proxy.Register) *APILoader {
	return &APILoader{register: register}
}

// RegisterAPIs load application middleware
func (m *APILoader) RegisterAPIs(cfgs []*api.Definition) {
	for _, spec := range cfgs {
		m.RegisterAPI(spec)
	}
}

// RegisterAPI register an API Definition in the register
func (m *APILoader) RegisterAPI(def *api.Definition) {
	logger := log.WithField("api_name", def.Name)
	logger.Debug("Starting RegisterAPI")

	active, err := def.Validate()
	if false == active && err != nil {
		logger.WithError(err).Error("Validation errors")
	}

	if false == def.Active {
		logger.Warn("API is not active, skipping...")
		active = false
	}

	if active {
		route := proxy.NewRoute(def.Proxy)

		for _, plg := range def.Plugins {
			l := logger.WithField("name", plg.Name)
			if plg.Enabled {
				l.Debug("Plugin enabled")

				setup, err := plugin.DirectiveAction(plg.Name)
				if err != nil {
					l.WithError(err).Error("Error loading plugin")
					continue
				}

				err = setup(def, route, plg.Config)
				if err != nil {
					l.WithError(err).Error("Error executing plugin")
				}
			} else {
				l.Debug("Plugin not enabled")
			}
		}

		if len(def.Proxy.Hosts) > 0 {
			route.AddInbound(middleware.NewHostMatcher(def.Proxy.Hosts).Handler)
		}

		m.register.Add(route)
		logger.Debug("API registered")
	} else {
		logger.WithError(err).Warn("API URI is invalid or not active, skipping...")
	}
}
