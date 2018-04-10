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
	configs  []*api.Spec
}

// NewAPILoader creates a new instance of the api manager
func NewAPILoader(register *proxy.Register) *APILoader {
	return &APILoader{register: register}
}

// RegisterAPIs load application middleware
func (m *APILoader) RegisterAPIs(cfgs []*api.Spec) {
	m.configs = cfgs

	for _, spec := range m.configs {
		m.RegisterAPI(spec)
	}
}

// RegisterAPI register an API Spec in the register
func (m *APILoader) RegisterAPI(referenceSpec *api.Spec) {
	logger := log.WithField("api_name", referenceSpec.Name)
	logger.Debug("Starting RegisterAPI")

	active, err := referenceSpec.Validate()
	if false == active && err != nil {
		logger.WithError(err).Error("Validation errors")
	}

	if false == referenceSpec.Active {
		logger.Warn("API is not active, skipping...")
		active = false
	}

	if active {
		route := proxy.NewRoute(referenceSpec.Proxy)

		for _, pDefinition := range referenceSpec.Plugins {
			l := logger.WithField("name", pDefinition.Name)
			if pDefinition.Enabled {
				l.Debug("Plugin enabled")

				setup, err := plugin.DirectiveAction(pDefinition.Name)
				if err != nil {
					l.WithError(err).Error("Error loading plugin")
					continue
				}

				err = setup(route, pDefinition.Config)
				if err != nil {
					l.WithError(err).Error("Error executing plugin")
				}
			} else {
				l.Debug("Plugin not enabled")
			}
		}

		if len(referenceSpec.Definition.Proxy.Hosts) > 0 {
			route.AddInbound(middleware.NewHostMatcher(referenceSpec.Definition.Proxy.Hosts).Handler)
		}

		m.register.Add(route)
		logger.Debug("API registered")
	} else {
		logger.WithError(err).Warn("API URI is invalid or not active, skipping...")
	}
}
