package janus

import (
	"encoding/json"
	"errors"

	"net/http"

	log "github.com/Sirupsen/logrus"
)

// OAuthMiddleware is the after middleware for OAuth routes
type OAuthMiddleware struct {
	oauthManager *OAuthManager
	oauthSpec    *OAuthSpec
}

// ProcessRequest is the middleware method.
func (m *OAuthMiddleware) ProcessRequest(rw http.ResponseWriter, req *http.Request) (int, error) {
	var newSession SessionState
	newSession.OAuthServerID = m.oauthSpec.ID

	log.WithFields(log.Fields{
		"req": req,
	}).Info("Getting body")

	var body []byte
	body = req.Context().Value("body").([]byte)

	if nil == body {
		return http.StatusInternalServerError, errors.New("Request body not present")
	}

	if marshalErr := json.Unmarshal(body, &newSession); marshalErr != nil {
		return http.StatusInternalServerError, marshalErr
	}

	m.oauthManager.Set(newSession.AccessToken, newSession, newSession.ExpiresIn)

	return http.StatusOK, nil
}
