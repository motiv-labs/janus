package main

import log "github.com/Sirupsen/logrus"

type OAuthRegister struct{}

func (m *OAuthRegister) GetProxiesForServer(oauth *OAuth) []Proxy {
	log.Debug("Loading oauth configuration")
	var proxies []Proxy

	//oauth proxy
	log.Debug("Registering authorize endpoint")
	authorizeProxy := oauth.OauthEndpoints.Authorize
	if validateProxy(authorizeProxy) {
		proxies = append(proxies, authorizeProxy)
	} else {
		log.Debug("No authorize endpoint")
	}

	log.Debug("Registering token endpoint")
	tokenProxy := oauth.OauthEndpoints.Token
	if validateProxy(tokenProxy) {
		proxies = append(proxies, tokenProxy)
	} else {
		log.Debug("No token endpoint")
	}

	log.Debug("Registering info endpoint")
	infoProxy := oauth.OauthEndpoints.Info
	if validateProxy(infoProxy) {
		proxies = append(proxies, infoProxy)
	} else {
		log.Debug("No info endpoint")
	}

	log.Debug("Registering create client endpoint")
	createProxy := oauth.OauthClientEndpoints.Create
	if validateProxy(createProxy) {
		proxies = append(proxies, createProxy)
	} else {
		log.Debug("No client create endpoint")
	}

	log.Debug("Registering remove client endpoint")
	removeProxy := oauth.OauthClientEndpoints.Remove
	if validateProxy(removeProxy) {
		proxies = append(proxies, removeProxy)
	} else {
		log.Debug("No client remove endpoint")
	}

	return proxies
}
