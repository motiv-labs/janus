package janus

import (
	"errors"
	"strings"

	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

type Oauth2KeyExists struct {
	*Middleware
	OAuthManager *OAuthManager
}

func (m *Oauth2KeyExists) ProcessRequest(req *http.Request, c *gin.Context) (error, int) {
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

		return errors.New("Authorization field missing"), http.StatusBadRequest
	}

	if strings.ToLower(parts[0]) != "bearer" {
		logger.Info("Bearer token malformed")
		return errors.New("Bearer token malformed"), http.StatusBadRequest
	}

	accessToken := parts[1]
	thisSessionState, keyExists := m.CheckSessionAndIdentityForValidKey(accessToken)

	if !keyExists || thisSessionState.OAuthServerID != m.Spec.OAuthServerID {
		log.WithFields(log.Fields{
			"path":   req.RequestURI,
			"origin": req.RemoteAddr,
			"key":    accessToken,
		}).Info("Attempted access with non-existent key.")

		return errors.New("Key not authorised"), http.StatusUnauthorized
	}

	c.Set(SessionData, thisSessionState)
	c.Set(AuthHeaderValue, accessToken)

	return nil, http.StatusOK
}

func (o *Oauth2KeyExists) CheckSessionAndIdentityForValidKey(key string) (SessionState, bool) {
	var thisSession SessionState

	//Checks if the key is present on the cache and if it didn't expire yet
	log.Debug("Querying keystore")
	if !o.OAuthManager.KeyExists(key) {
		log.Debug("Key not found in keystore")
		return thisSession, false
	}

	// 2. If not there, get it from the AuthorizationHandler
	return o.OAuthManager.IsKeyAuthorised(key)
}
