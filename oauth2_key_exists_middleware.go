package janus

import (
	"context"
	"errors"
	"strings"

	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/hellofresh/janus/response"
)

// Oauth2KeyExistsMiddleware checks the integrity of the provided OAuth headers
type Oauth2KeyExistsMiddleware struct {
	*Middleware
	OAuthManager *OAuthManager
}

// Serve is the middleware method.
func (m *Oauth2KeyExistsMiddleware) Serve(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Starting Oauth2KeyExists middleware")
		logger := log.WithFields(log.Fields{
			"path":   r.RequestURI,
			"origin": r.RemoteAddr,
		})

		// We're using OAuth, start checking for access keys
		authHeaderValue := r.Header.Get("Authorization")
		parts := strings.Split(authHeaderValue, " ")
		if len(parts) < 2 {
			logger.Info("Attempted access with malformed header, no auth header found.")
			response.JSON(w, http.StatusBadRequest, errors.New("authorization field missing"))
		}

		if strings.ToLower(parts[0]) != "bearer" {
			logger.Info("Bearer token malformed")
			response.JSON(w, http.StatusBadRequest, errors.New("bearer token malformed"))
		}

		accessToken := parts[1]
		thisSessionState, keyExists := m.CheckSessionAndIdentityForValidKey(accessToken)

		if !keyExists || thisSessionState.OAuthServerID != m.Spec.OAuthServerID {
			log.WithFields(log.Fields{
				"path":   r.RequestURI,
				"origin": r.RemoteAddr,
				"key":    accessToken,
			}).Info("Attempted access with non-existent key.")

			response.JSON(w, http.StatusUnauthorized, errors.New("key not authorised"))
			return
		}

		context.WithValue(r.Context(), SessionData, thisSessionState)
		context.WithValue(r.Context(), AuthHeaderValue, accessToken)

		handler.ServeHTTP(w, r)
	})
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
