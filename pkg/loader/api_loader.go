package loader

import (
	"github.com/afex/hystrix-go/hystrix"
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
	return &APILoader{register}
}

// LoadDefinitions registers all ApiDefinitions from a data source
func (m *APILoader) LoadDefinitions(repo api.Repository) {
	specs := m.getAPISpecs(repo)
	m.RegisterApis(specs)
}

// RegisterApis load application middleware
func (m *APILoader) RegisterApis(apiSpecs []*api.Spec) {
	for _, referenceSpec := range apiSpecs {
		m.RegisterAPI(referenceSpec)
	}
}

// RegisterAPI register an API Spec in the register
func (m *APILoader) RegisterAPI(referenceSpec *api.Spec) {
	logger := log.WithField("api_name", referenceSpec.Name)

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

// getAPISpecs Load application specs from data source
func (m *APILoader) getAPISpecs(repo api.Repository) []*api.Spec {
	definitions, err := repo.FindAll()
	if err != nil {
		log.Panic(err)
	}

	var specs []*api.Spec
	for _, d := range definitions {
		m.createCircuitBreakerDefinition(d)
		specs = append(specs, &api.Spec{Definition: d})
	}

	return specs
}

func (m *APILoader) createCircuitBreakerDefinition(d *api.Definition) {
	hystrix.ConfigureCommand(d.Proxy.ListenPath, hystrix.CommandConfig{
		Timeout:                d.CircuitBreaker.Timeout,
		MaxConcurrentRequests:  d.CircuitBreaker.MaxConcurrentRequests,
		ErrorPercentThreshold:  d.CircuitBreaker.ErrorPercentThreshold,
		RequestVolumeThreshold: d.CircuitBreaker.RequestVolumeThreshold,
		SleepWindow:            d.CircuitBreaker.SleepWindow,
	})
}
