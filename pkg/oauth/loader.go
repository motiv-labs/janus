package oauth

import (
	log "github.com/Sirupsen/logrus"
	"github.com/hellofresh/janus/pkg/cors"
	"github.com/hellofresh/janus/pkg/middleware"
	"github.com/hellofresh/janus/pkg/proxy"
	"github.com/hellofresh/janus/pkg/store"
)

// Loader handles the loading of the api specs
type Loader struct {
	register *proxy.Register
	store    store.Store
	accessor *middleware.DatabaseAccessor
	debug    bool
}

// NewLoader creates a new instance of the api manager
func NewLoader(register *proxy.Register, store store.Store, accessor *middleware.DatabaseAccessor, debug bool) *Loader {
	return &Loader{register, store, accessor, debug}
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
		log.Debug("Loading oauth configuration")
		corsHandler := cors.NewMiddleware(oauthServer.CorsMeta, m.debug).Handler
		//oauth proxy
		log.Debug("Registering authorize endpoint")
		authorizeProxy := oauthServer.Endpoints.Authorize
		if proxy.Validate(authorizeProxy) {
			m.register.Add(proxy.NewRoute(authorizeProxy, corsHandler))
		} else {
			log.Debug("No authorize endpoint")
		}

		log.Debug("Registering token endpoint")
		tokenProxy := oauthServer.Endpoints.Token
		if proxy.Validate(tokenProxy) {
			m.register.Add(proxy.NewRoute(tokenProxy, NewSecretMiddleware(oauthServer).Handler, corsHandler))
		} else {
			log.Debug("No token endpoint")
		}

		log.Debug("Registering info endpoint")
		infoProxy := oauthServer.Endpoints.Info
		if proxy.Validate(infoProxy) {
			m.register.Add(proxy.NewRoute(infoProxy, corsHandler))
		} else {
			log.Debug("No info endpoint")
		}

		log.Debug("Registering revoke endpoint")
		revokeProxy := oauthServer.Endpoints.Revoke
		if proxy.Validate(revokeProxy) {
			m.register.Add(proxy.NewRoute(revokeProxy, corsHandler, NewRevokeMiddleware(m.store).Handler))
		} else {
			log.Debug("No revoke endpoint")
		}

		log.Debug("Registering create client endpoint")
		createProxy := oauthServer.ClientEndpoints.Create
		if proxy.Validate(createProxy) {
			m.register.Add(proxy.NewRoute(createProxy, corsHandler))
		} else {
			log.Debug("No client create endpoint")
		}

		log.Debug("Registering remove client endpoint")
		removeProxy := oauthServer.ClientEndpoints.Remove
		if proxy.Validate(removeProxy) {
			m.register.Add(proxy.NewRoute(removeProxy, corsHandler))
		} else {
			log.Debug("No client remove endpoint")
		}
	}

	log.Debug("Done loading OAuth servers configurations")
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
