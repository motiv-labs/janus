package oauth

import (
	log "github.com/Sirupsen/logrus"
	"github.com/hellofresh/janus/pkg/proxy"
	"github.com/hellofresh/janus/pkg/store"
	"github.com/rs/cors"
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
		corsHandler := cors.New(cors.Options{
			AllowedOrigins:   oauthServer.CorsMeta.Domains,
			AllowedMethods:   oauthServer.CorsMeta.Methods,
			AllowedHeaders:   oauthServer.CorsMeta.RequestHeaders,
			ExposedHeaders:   oauthServer.CorsMeta.ExposedHeaders,
			AllowCredentials: true,
		}).Handler

		log.Debug("Registering authorize endpoint")
		authorizeProxy := oauthServer.Endpoints.Authorize
		if isValid, err := authorizeProxy.Validate(); isValid && err == nil {
			m.register.Add(proxy.NewRoute(authorizeProxy, corsHandler))
		} else {
			log.WithError(err).Debug("No authorize endpoint")
		}

		log.Debug("Registering token endpoint")
		tokenProxy := oauthServer.Endpoints.Token
		if isValid, err := tokenProxy.Validate(); isValid && err == nil {
			m.register.AddWithInOut(
				proxy.NewRoute(tokenProxy, NewSecretMiddleware(oauthServer).Handler, corsHandler),
				nil,
				proxy.NewOutChain(NewTokenPlugin(m.storage, repo).Out),
			)
		} else {
			log.WithError(err).Debug("No token endpoint")
		}

		log.Debug("Registering info endpoint")
		infoProxy := oauthServer.Endpoints.Info
		if isValid, err := infoProxy.Validate(); isValid && err == nil {
			m.register.Add(proxy.NewRoute(infoProxy, corsHandler))
		} else {
			log.WithError(err).Debug("No info endpoint")
		}

		log.Debug("Registering revoke endpoint")
		revokeProxy := oauthServer.Endpoints.Revoke
		if isValid, err := revokeProxy.Validate(); isValid && err == nil {
			m.register.Add(proxy.NewRoute(revokeProxy, corsHandler, NewRevokeMiddleware(oauthServer).Handler))
		} else {
			log.WithError(err).Debug("No revoke endpoint")
		}

		log.Debug("Registering create client endpoint")
		createProxy := oauthServer.ClientEndpoints.Create
		if isValid, err := createProxy.Validate(); isValid && err == nil {
			m.register.Add(proxy.NewRoute(createProxy, corsHandler))
		} else {
			log.WithError(err).Debug("No client create endpoint")
		}

		log.Debug("Registering remove client endpoint")
		removeProxy := oauthServer.ClientEndpoints.Remove
		if isValid, err := createProxy.Validate(); isValid && err == nil {
			m.register.Add(proxy.NewRoute(removeProxy, corsHandler))
		} else {
			log.WithError(err).Debug("No client remove endpoint")
		}
	}

	log.Debug("Done loading OAuth servers configurations")
}

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
