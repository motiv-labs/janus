package janus

import (
	"encoding/json"
	"errors"

	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/hellofresh/janus/oauth"
	"github.com/hellofresh/janus/response"
	"github.com/hellofresh/janus/session"
)

// OAuthMiddleware is the after middleware for OAuth routes
type OAuthMiddleware struct {
	manager   *oauth.Manager
	oauthSpec *OAuthSpec
}

// Serve is the middleware method.
func (m *OAuthMiddleware) Serve(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var newSession session.SessionState
		newSession.OAuthServerID = m.oauthSpec.ID

		log.WithFields(log.Fields{
			"req": r,
		}).Info("Getting body")

		body := r.Context().Value("body")
		if nil == body {
			response.JSON(w, http.StatusInternalServerError, errors.New("request body not present"))
			return
		}

		if marshalErr := json.Unmarshal(body.([]byte), &newSession); marshalErr != nil {
			response.JSON(w, http.StatusInternalServerError, marshalErr)
			return
		}

		m.manager.Set(newSession.AccessToken, newSession, newSession.ExpiresIn)

		handler.ServeHTTP(w, r)
	})
}
