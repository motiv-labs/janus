package oauth

import (
	log "github.com/Sirupsen/logrus"
	"github.com/hellofresh/janus/cors"
	"github.com/hellofresh/janus/middleware"
	"github.com/hellofresh/janus/proxy"
	"github.com/hellofresh/janus/router"
)

// Loader handles the loading of the api specs
type Loader struct {
	proxyRegister *proxy.Register
	accessor      *middleware.DatabaseAccessor
	debug         bool
}

// NewLoader creates a new instance of the api manager
func NewLoader(router router.Router, accessor *middleware.DatabaseAccessor, proxyRegister *proxy.Register, debug bool) *Loader {
	return &Loader{proxyRegister, accessor, debug}
}

// Load loads all api specs from a datasource
func (m *Loader) Load() {
	oAuthServers := m.getOAuthServers()
	go m.LoadOAuthServers(oAuthServers)
}

// LoadOAuthServers loads and register the oauth servers
func (m *Loader) LoadOAuthServers(oauthServers []*OAuth) {
	log.Debug("Loading OAuth servers configurations")

	for _, oauthServer := range oauthServers {
		proxies := GetProxiesForServer(oauthServer)
		m.proxyRegister.RegisterMany(
			proxies,
			NewSecretMiddleware(oauthServer).Handler,
			cors.NewMiddleware(oauthServer.CorsMeta, m.debug).Handler,
		)
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
