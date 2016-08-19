package main

import (
	"errors"
	"strings"

	log "github.com/Sirupsen/logrus"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/context"
)

type Oauth2KeyExists struct {
	*Middleware
}

//Important staff, iris middleware must implement the iris.Handler interface which is:
func (m Oauth2KeyExists) ProcessRequest(req *http.Request, c *gin.Context) (error, int) {
	m.Logger.Debug("Starting Oauth2KeyExists middleware")

	if false == m.Spec.UseOauth2 {
		m.Logger.Debug("OAuth2 not enabled")
		return nil, http.StatusOK
	}

	// We're using OAuth, start checking for access keys
	authHeaderValue := string(req.Header.Get("Authorization"))
	parts := strings.Split(authHeaderValue, " ")
	if len(parts) < 2 {
		m.Logger.Info("Attempted access with malformed header, no auth header found.")

		return errors.New("Authorization field missing"), http.StatusBadRequest
	}

	if strings.ToLower(parts[0]) != "bearer" {
		m.Logger.Info("Bearer token malformed")

		return errors.New("Bearer token malformed"), http.StatusBadRequest
	}

	accessToken := parts[1]
	thisSessionState, keyExists := m.CheckSessionAndIdentityForValidKey(accessToken)

	if !keyExists {
		m.Logger.WithFields(log.Fields{
			"path":  req.RequestURI,
			"origin": req.RemoteAddr,
			"key":    accessToken,
		}).Info("Attempted access with non-existent key.")

		return errors.New("Key not authorised"), http.StatusUnauthorized
	}

	context.Set(req, SessionData, thisSessionState)
	context.Set(req, AuthHeaderValue, accessToken)

	return nil, http.StatusOK
}
