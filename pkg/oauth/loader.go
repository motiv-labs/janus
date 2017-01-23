package oauth

import (
	log "github.com/Sirupsen/logrus"
	"github.com/hellofresh/janus/pkg/cors"
	"github.com/hellofresh/janus/pkg/middleware"
	"github.com/hellofresh/janus/pkg/proxy"
)

// Loader handles the loading of the api specs
type Loader struct {
	register *proxy.Register
	accessor *middleware.DatabaseAccessor
	debug    bool
}

// NewLoader creates a new instance of the api manager
func NewLoader(register *proxy.Register, accessor *middleware.DatabaseAccessor, debug bool) *Loader {
	return &Loader{register, accessor, debug}
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
		m.register.AddMany(m.RegisterOAuthServer(oauthServer))
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
