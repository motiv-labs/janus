package oauth

import (
	log "github.com/Sirupsen/logrus"
	"github.com/hellofresh/janus/proxy"
	"github.com/hellofresh/janus/router"
)

// GetRoutesForServer converts an oauth definition into many proxies
func GetRoutesForServer(oauth *OAuth, handlers ...router.Constructor) []*proxy.Route {
	log.Debug("Loading oauth configuration")
	var routes []*proxy.Route

	//oauth proxy
	log.Debug("Registering authorize endpoint")
	authorizeProxy := oauth.Endpoints.Authorize
	if proxy.Validate(authorizeProxy) {
		routes = append(routes, proxy.NewRoute(authorizeProxy, handlers...))
	} else {
		log.Debug("No authorize endpoint")
	}

	log.Debug("Registering token endpoint")
	tokenProxy := oauth.Endpoints.Token
	if proxy.Validate(tokenProxy) {
		routes = append(routes, proxy.NewRoute(tokenProxy, handlers...))
	} else {
		log.Debug("No token endpoint")
	}

	log.Debug("Registering info endpoint")
	infoProxy := oauth.Endpoints.Info
	if proxy.Validate(infoProxy) {
		routes = append(routes, proxy.NewRoute(infoProxy, handlers...))
	} else {
		log.Debug("No info endpoint")
	}

	log.Debug("Registering create client endpoint")
	createProxy := oauth.ClientEndpoints.Create
	if proxy.Validate(createProxy) {
		routes = append(routes, proxy.NewRoute(createProxy, handlers...))
	} else {
		log.Debug("No client create endpoint")
	}

	log.Debug("Registering remove client endpoint")
	removeProxy := oauth.ClientEndpoints.Remove
	if proxy.Validate(removeProxy) {
		routes = append(routes, proxy.NewRoute(removeProxy, handlers...))
	} else {
		log.Debug("No client remove endpoint")
	}

	return routes
}
