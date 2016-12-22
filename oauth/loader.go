package oauth

import (
	"github.com/hellofresh/janus/cors"
	"github.com/hellofresh/janus/loader"
	"github.com/hellofresh/janus/log"
	"github.com/hellofresh/janus/middleware"
	"github.com/hellofresh/janus/proxy"
)

// Loader handles the loading of the api specs
type Loader struct {
	registerChan *proxy.RegisterChan
	accessor     *middleware.DatabaseAccessor
	debug        bool
}

// NewLoader creates a new instance of the api manager
func NewLoader(registerChan *proxy.RegisterChan, accessor *middleware.DatabaseAccessor, debug bool) *Loader {
	return &Loader{registerChan, accessor, debug}
}

// ListenToChanges listens to any changes that might require a reload of configurations
func (m *Loader) ListenToChanges(tracker *loader.Tracker) {
	go func() {
		for {
			select {
			case <-tracker.StopChan():
				log.Debug("Stopping listening to api changes....")
				return
			case <-tracker.Changed():
				log.Debug("A change was identified. Reloading api configurations....")
				m.Load()
			}
		}
	}()
}

// Load loads all api specs from a datasource
func (m *Loader) Load() {
	oAuthServers := m.getOAuthServers()
	m.RegisterOAuthServers(oAuthServers)
}

// RegisterOAuthServers register many oauth servers
func (m *Loader) RegisterOAuthServers(oauthServers []*OAuth) {
	log.Debug("Loading OAuth servers configurations")

	for _, oauthServer := range oauthServers {
		m.registerChan.Many <- m.RegisterOAuthServer(oauthServer)
	}

	log.Debug("Done loading OAuth servers configurations")
}

// RegisterOAuthServer register the an oauth server
func (m *Loader) RegisterOAuthServer(oauthServer *OAuth) []*proxy.Route {
	return GetRoutesForServer(
		oauthServer,
		NewSecretMiddleware(oauthServer).Handler,
		cors.NewMiddleware(oauthServer.CorsMeta, m.debug).Handler,
	)
}

//getOAuthServers Load oauth servers from datasource
func (m *Loader) getOAuthServers() []*OAuth {
	log.Debug("Using Oauth servers configuration from Mongo DB")
	repo, err := NewMongoRepository(m.accessor.Session.DB(""))
	if err != nil {
		log.Panic(err)
	}

	oauthServers, err := repo.FindAll()
	if err != nil {
		log.Panic(err)
	}

	return oauthServers
}
