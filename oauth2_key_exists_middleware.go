package janus

import (
	"context"
	"errors"
	"strings"

	"net/http"

	log "github.com/Sirupsen/logrus"
)

// Oauth2KeyExistsMiddleware checks the integrity of the provided OAuth headers
type Oauth2KeyExistsMiddleware struct {
	*Middleware
	OAuthManager *OAuthManager
}

// ProcessRequest is the middleware method.
func (m *Oauth2KeyExistsMiddleware) ProcessRequest(rw http.ResponseWriter, req *http.Request) (int, error) {
	log.Debug("Starting Oauth2KeyExists middleware")
	logger := log.WithFields(log.Fields{
		"path":   req.RequestURI,
		"origin": req.RemoteAddr,
	})

	// We're using OAuth, start checking for access keys
	authHeaderValue := req.Header.Get("Authorization")
	parts := strings.Split(authHeaderValue, " ")
	if len(parts) < 2 {
		logger.Info("Attempted access with malformed header, no auth header found.")

		return http.StatusBadRequest, errors.New("Authorization field missing")
	}

	if strings.ToLower(parts[0]) != "bearer" {
		logger.Info("Bearer token malformed")
		return http.StatusBadRequest, errors.New("Bearer token malformed")
	}

	accessToken := parts[1]
	thisSessionState, keyExists := m.CheckSessionAndIdentityForValidKey(accessToken)

	if !keyExists || thisSessionState.OAuthServerID != m.Spec.OAuthServerID {
		log.WithFields(log.Fields{
			"path":   req.RequestURI,
			"origin": req.RemoteAddr,
			"key":    accessToken,
		}).Info("Attempted access with non-existent key.")

		return http.StatusUnauthorized, errors.New("Key not authorised")
	}

	context.WithValue(req.Context(), SessionData, thisSessionState)
	context.WithValue(req.Context(), AuthHeaderValue, accessToken)

	return http.StatusOK, nil
}

// CheckSessionAndIdentityForValidKey ensures we have the valid key in the session store
func (m *Oauth2KeyExistsMiddleware) CheckSessionAndIdentityForValidKey(key string) (SessionState, bool) {
	var thisSession SessionState

	// Checks if the key is present on the cache and if it didn't expire yet
	log.Debug("Querying keystore")
	if !m.OAuthManager.KeyExists(key) {
		log.Debug("Key not found in keystore")
		return thisSession, false
	}

	// 2. If not there, get it from the AuthorizationHandler
	return m.OAuthManager.IsKeyAuthorised(key)
}
