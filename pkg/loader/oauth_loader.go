package loader

import (
	log "github.com/sirupsen/logrus"
	"github.com/hellofresh/janus/pkg/oauth"
	"github.com/hellofresh/janus/pkg/proxy"
	"github.com/hellofresh/janus/pkg/store"
	"github.com/rs/cors"
)

// OAuthLoader handles the loading of the api specs
type OAuthLoader struct {
	register *proxy.Register
	storage  store.Store
}

// NewOAuthLoader creates a new instance of the Loader
func NewOAuthLoader(register *proxy.Register, storage store.Store) *OAuthLoader {
	return &OAuthLoader{register, storage}
}

// LoadDefinitions loads all oauth servers from a data source
func (m *OAuthLoader) LoadDefinitions(repo oauth.Repository) {
	oAuthServers := m.getOAuthServers(repo)
	m.RegisterOAuthServers(oAuthServers, repo)
}

// RegisterOAuthServers register many oauth servers
func (m *OAuthLoader) RegisterOAuthServers(oauthServers []*oauth.Spec, repo oauth.Repository) {
	log.Debug("Loading OAuth servers configurations")

	for _, oauthServer := range oauthServers {
		logger := log.WithField("name", oauthServer.Name)
		logger.Debug("Registering OAuth server")

		corsHandler := cors.New(cors.Options{
			AllowedOrigins:   oauthServer.CorsMeta.Domains,
			AllowedMethods:   oauthServer.CorsMeta.Methods,
			AllowedHeaders:   oauthServer.CorsMeta.RequestHeaders,
			ExposedHeaders:   oauthServer.CorsMeta.ExposedHeaders,
			AllowCredentials: true,
		}).Handler

		logger.Debug("Registering authorize endpoint")
		authorizeProxy := oauthServer.Endpoints.Authorize
		if isValid, err := authorizeProxy.Validate(); isValid && err == nil {
			m.register.Add(proxy.NewRoute(authorizeProxy, corsHandler))
		} else {
			logger.WithError(err).Debug("No authorize endpoint")
		}

		logger.Debug("Registering token endpoint")
		tokenProxy := oauthServer.Endpoints.Token
		if isValid, err := tokenProxy.Validate(); isValid && err == nil {
			m.register.AddWithInOut(
				proxy.NewRoute(tokenProxy, oauth.NewSecretMiddleware(oauthServer).Handler, corsHandler),
				nil,
				proxy.NewOutChain(oauth.NewTokenPlugin(m.storage, repo).Out),
			)
		} else {
			logger.WithError(err).Debug("No token endpoint")
		}

		logger.Debug("Registering info endpoint")
		infoProxy := oauthServer.Endpoints.Info
		if isValid, err := infoProxy.Validate(); isValid && err == nil {
			m.register.Add(proxy.NewRoute(infoProxy, corsHandler))
		} else {
			logger.WithError(err).Debug("No info endpoint")
		}

		logger.Debug("Registering revoke endpoint")
		revokeProxy := oauthServer.Endpoints.Revoke
		if isValid, err := revokeProxy.Validate(); isValid && err == nil {
			m.register.Add(proxy.NewRoute(revokeProxy, corsHandler, oauth.NewRevokeMiddleware(oauthServer).Handler))
		} else {
			logger.WithError(err).Debug("No revoke endpoint")
		}

		logger.Debug("Registering create client endpoint")
		createProxy := oauthServer.ClientEndpoints.Create
		if isValid, err := createProxy.Validate(); isValid && err == nil {
			m.register.Add(proxy.NewRoute(createProxy, corsHandler))
		} else {
			logger.WithError(err).Debug("No client create endpoint")
		}

		logger.Debug("Registering remove client endpoint")
		removeProxy := oauthServer.ClientEndpoints.Remove
		if isValid, err := createProxy.Validate(); isValid && err == nil {
			m.register.Add(proxy.NewRoute(removeProxy, corsHandler))
		} else {
			logger.WithError(err).Debug("No client remove endpoint")
		}

		logger.Debug("Oauth server registered")
	}

	log.Debug("Done loading OAuth servers configurations")
}

func (m *OAuthLoader) getOAuthServers(repo oauth.Repository) []*oauth.Spec {
	oauthServers, err := repo.FindAll()
	if err != nil {
		log.Panic(err)
	}

	var specs []*oauth.Spec
	for _, oauthServer := range oauthServers {
		spec := new(oauth.Spec)
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

func (m *OAuthLoader) getManager(oauthServer *oauth.OAuth) (oauth.Manager, error) {
	managerType, err := oauth.ParseType(oauthServer.TokenStrategy.Name)
	if nil != err {
		return nil, err
	}

	return oauth.NewManagerFactory(m.storage, oauthServer.TokenStrategy.Settings).Build(managerType)
}
