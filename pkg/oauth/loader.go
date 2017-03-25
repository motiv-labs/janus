package oauth

import (
	log "github.com/Sirupsen/logrus"
	"github.com/hellofresh/janus/pkg/cors"
	"github.com/hellofresh/janus/pkg/proxy"
	"github.com/hellofresh/janus/pkg/store"
)

// Loader handles the loading of the api specs
type Loader struct {
	register *proxy.Register
	storage  store.Store
}

// NewLoader creates a new instance of the Loader
func NewLoader(register *proxy.Register, storage store.Store) *Loader {
	return &Loader{register, storage}
}

// LoadDefinitions loads all oauth servers from a datasource
func (m *Loader) LoadDefinitions(repo Repository) {
	oAuthServers := m.getOAuthServers(repo)
	m.RegisterOAuthServers(oAuthServers, repo)
}

// RegisterOAuthServers register many oauth servers
func (m *Loader) RegisterOAuthServers(oauthServers []*Spec, repo Repository) {
	log.Debug("Loading OAuth servers configurations")

	for _, oauthServer := range oauthServers {
		log.Debug("Loading oauth configuration")
		corsHandler := cors.NewMiddleware(oauthServer.CorsMeta, false).Handler
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
			m.register.AddWithInOut(
				proxy.NewRoute(tokenProxy, NewSecretMiddleware(oauthServer).Handler, corsHandler),
				nil,
				proxy.NewOutChain(NewTokenPlugin(m.storage, repo).Out),
			)
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
			m.register.Add(proxy.NewRoute(revokeProxy, corsHandler, NewRevokeMiddleware(oauthServer).Handler))
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
func (m *Loader) getOAuthServers(repo Repository) []*Spec {
	oauthServers, err := repo.FindAll()
	if err != nil {
		log.Panic(err)
	}

	var specs []*Spec
	for _, oauthServer := range oauthServers {
		spec := new(Spec)
		spec.OAuth = oauthServer
		manager, err := m.getManager(oauthServer)
		if nil != err {
			log.WithError(err).Error("Oauth definition is not well configured, skipping...")
			continue
		}
		spec.Manager = manager
		specs = append(specs, spec)
	}

	return specs
}

func (m *Loader) getManager(oauthServer *OAuth) (Manager, error) {
	managerType, err := ParseType(oauthServer.TokenStrategy.Name)
	if nil != err {
		return nil, err
	}

	return NewManagerFactory(m.storage, oauthServer.TokenStrategy.Settings).Build(managerType)
}
