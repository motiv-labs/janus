package oauth

import (
	log "github.com/Sirupsen/logrus"
	"github.com/hellofresh/janus/proxy"
)

// GetProxiesForServer converts an oauth definition into many proxies
func GetProxiesForServer(oauth *OAuth) []proxy.Proxy {
	log.Debug("Loading oauth configuration")
	var proxies []proxy.Proxy

	//oauth proxy
	log.Debug("Registering authorize endpoint")
	authorizeProxy := oauth.Endpoints.Authorize
	if proxy.Validate(authorizeProxy) {
		proxies = append(proxies, authorizeProxy)
	} else {
		log.Debug("No authorize endpoint")
	}

	log.Debug("Registering token endpoint")
	tokenProxy := oauth.Endpoints.Token
	if proxy.Validate(tokenProxy) {
		proxies = append(proxies, tokenProxy)
	} else {
		log.Debug("No token endpoint")
	}

	log.Debug("Registering info endpoint")
	infoProxy := oauth.Endpoints.Info
	if proxy.Validate(infoProxy) {
		proxies = append(proxies, infoProxy)
	} else {
		log.Debug("No info endpoint")
	}

	log.Debug("Registering create client endpoint")
	createProxy := oauth.ClientEndpoints.Create
	if proxy.Validate(createProxy) {
		proxies = append(proxies, createProxy)
	} else {
		log.Debug("No client create endpoint")
	}

	log.Debug("Registering remove client endpoint")
	removeProxy := oauth.ClientEndpoints.Remove
	if proxy.Validate(removeProxy) {
		proxies = append(proxies, removeProxy)
	} else {
		log.Debug("No client remove endpoint")
	}

	return proxies
}
