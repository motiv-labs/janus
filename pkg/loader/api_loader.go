package loader

import (
	log "github.com/Sirupsen/logrus"
	"github.com/hellofresh/janus/pkg/api"
	"github.com/hellofresh/janus/pkg/plugin"
	"github.com/hellofresh/janus/pkg/proxy"
	"github.com/hellofresh/janus/pkg/router"
)

// Loader is responsible for loading all apis form a datastore and configure them in a register
type Loader struct {
	register     *proxy.Register
	pluginLoader *plugin.Loader
}

// NewLoader creates a new instance of the api manager
func NewLoader(register *proxy.Register, pluginLoader *plugin.Loader) *Loader {
	return &Loader{register, pluginLoader}
}

// LoadDefinitions will connect and download ApiDefintions from a Mongo DB instance.
func (m *Loader) LoadDefinitions(repo api.Repository) {
	specs := m.getAPISpecs(repo)
	m.RegisterApis(specs)
}

// RegisterApis load application middleware
func (m *Loader) RegisterApis(apiSpecs []*api.Spec) {
	for _, referenceSpec := range apiSpecs {
		m.RegisterAPI(referenceSpec)
	}
}

// RegisterAPI register an API Spec in the register
func (m *Loader) RegisterAPI(referenceSpec *api.Spec) {
	logger := log.WithField("api_name", referenceSpec.Name)

	active, err := referenceSpec.Validate()
	if false == active && err != nil {
		logger.WithField("errors", err.Error()).Warn("Validation errors")
	}

	if false == referenceSpec.Active {
		logger.Warn("API is not active, skiping...")
		active = false
	}

	if active {
		var handlers []router.Constructor

		for pName, pDefinition := range referenceSpec.Plugins {
			pDefinition.Name = pName
			if pDefinition.Enabled {
				logger.Debugf("Plugin %s enabled", pName)
				if p := m.pluginLoader.Get(pDefinition.Name); p != nil {
					middlewares, err := p.GetMiddlewares(pDefinition.Config, referenceSpec)
					if err != nil {
						logger.WithError(err).
							WithField("plugin_name", pName).
							Error("Error loading plugin")
					}

					for _, mw := range middlewares {
						handlers = append(handlers, mw)
					}
				}
			} else {
				logger.Debugf("Plugin %s not enabled", pName)
			}
		}

		m.register.Add(proxy.NewRoute(referenceSpec.Proxy, handlers...))
		logger.Debug("API registered")
	} else {
		logger.Warn("API URI is invalid or not active, skipping...")
	}
}

//getAPISpecs Load application specs from datasource
func (m *Loader) getAPISpecs(repo api.Repository) []*api.Spec {
	definitions, err := repo.FindAll()
	if err != nil {
		log.Panic(err)
	}

	var specs []*api.Spec
	for _, definition := range definitions {
		specs = append(specs, &api.Spec{Definition: definition})
	}

	return specs
}
