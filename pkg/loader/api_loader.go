package loader

import (
	"encoding/json"

	log "github.com/Sirupsen/logrus"
	"github.com/hellofresh/janus/pkg/api"
	"github.com/hellofresh/janus/pkg/middleware"
	"github.com/hellofresh/janus/pkg/plugin"
	"github.com/hellofresh/janus/pkg/proxy"
	"github.com/hellofresh/janus/pkg/router"
	"github.com/hellofresh/janus/pkg/store"
)

// APILoader is responsible for loading all apis form a datastore and configure them in a register
type APILoader struct {
	register     *proxy.Register
	pluginLoader *plugin.Loader
	subs         *store.Subscription
}

// NewAPILoader creates a new instance of the api manager
func NewAPILoader(register *proxy.Register, pluginLoader *plugin.Loader, subs *store.Subscription) *APILoader {
	return &APILoader{register, pluginLoader, subs}
}

// LoadDefinitions will connect and download ApiDefintions from a Mongo DB instance.
func (m *APILoader) LoadDefinitions(repo api.Repository) {
	specs := m.getAPISpecs(repo)
	m.RegisterApis(specs)
}

// RegisterApis load application middleware
func (m *APILoader) RegisterApis(apiSpecs []*api.Spec) {
	for _, referenceSpec := range apiSpecs {
		if m.subs != nil {
			log.Debug("Listening for changes on for the API definitions")
			go m.listenForChanges(referenceSpec.Definition)
		}
		m.RegisterAPI(referenceSpec)
	}
}

// RegisterAPI register an API Spec in the register
func (m *APILoader) RegisterAPI(referenceSpec *api.Spec) {
	logger := log.WithField("api_name", referenceSpec.Name)

	active, err := referenceSpec.Validate()
	if false == active && err != nil {
		logger.WithError(err).Warn("Validation errors")
	}

	if false == referenceSpec.Active {
		logger.Warn("API is not active, skiping...")
		active = false
	}

	if active {
		var handlers []router.Constructor

		for _, pDefinition := range referenceSpec.Plugins {
			if pDefinition.Enabled {
				logger.WithField("name", pDefinition.Name).Debug("Plugin enabled")
				if p := m.pluginLoader.Get(pDefinition.Name); p != nil {
					middlewares, err := p.GetMiddlewares(pDefinition.Config, referenceSpec)
					if err != nil {
						logger.WithError(err).
							WithField("plugin_name", pDefinition.Name).
							Error("Error loading plugin")
					}

					for _, mw := range middlewares {
						handlers = append(handlers, mw)
					}
				}
			} else {
				logger.WithField("name", pDefinition.Name).Debug("Plugin not enabled")
			}
		}

		if len(referenceSpec.Definition.Proxy.Hosts) > 0 {
			handlers = append(handlers, middleware.NewHostMatcher(referenceSpec.Definition.Proxy.Hosts).Handler)
		}

		m.register.Add(proxy.NewRoute(referenceSpec.Proxy, handlers...))
		logger.Debug("API registered")
	} else {
		logger.WithError(err).Warn("API URI is invalid or not active, skipping...")
	}
}

//getAPISpecs Load application specs from datasource
func (m *APILoader) getAPISpecs(repo api.Repository) []*api.Spec {
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

func (m *APILoader) listenForChanges(def *api.Definition) {
	for {
		select {
		case msg := <-m.subs.Message:
			var msgDefinition *api.Definition
			json.Unmarshal(msg, &msgDefinition)

			if def.Name == msgDefinition.Name {
				*def = *msgDefinition
			}
		}
	}
}
