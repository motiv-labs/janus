package main

import (
	"errors"
	"strings"

	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

type Oauth2KeyExists struct {
	*Middleware
}

func (m *Oauth2KeyExists) ProcessRequest(req *http.Request, c *gin.Context) (error, int) {
	log.Debug("Starting Oauth2KeyExists middleware")

	if false == m.Spec.UseOauth2 {
		log.Debug("OAuth2 not enabled")
		return nil, http.StatusOK
	}

	// We're using OAuth, start checking for access keys
	authHeaderValue := string(req.Header.Get("Authorization"))
	parts := strings.Split(authHeaderValue, " ")
	if len(parts) < 2 {
		log.Info("Attempted access with malformed header, no auth header found.")
		return errors.New("Authorization field missing"), http.StatusBadRequest
	}

	if strings.ToLower(parts[0]) != "bearer" {
		log.Info("Bearer token malformed")
		return errors.New("Bearer token malformed"), http.StatusBadRequest
	}

	accessToken := parts[1]
	thisSessionState, keyExists := m.CheckSessionAndIdentityForValidKey(accessToken)

	if !keyExists {
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
