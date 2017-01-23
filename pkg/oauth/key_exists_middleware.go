package oauth

import (
	"context"
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/hellofresh/janus/pkg/request"
	"github.com/hellofresh/janus/pkg/session"
	"gopkg.in/mgo.v2/bson"
)

// Enums for keys to be stored in a session context - this is how gorilla expects
// these to be implemented and is lifted pretty much from docs
var (
	SessionData     = request.ContextKey("session_data")
	AuthHeaderValue = request.ContextKey("auth_header")
)

// KeyExistsMiddleware checks the integrity of the provided OAuth headers
type KeyExistsMiddleware struct {
	manager *Manager
	// TODO pass api.Spec pointer here, currently impossible as circullar referens occurs
	oAuthServerID bson.ObjectId
}

// NewKeyExistsMiddleware creates a new instance of KeyExistsMiddleware
func NewKeyExistsMiddleware(manager *Manager, oAuthServerID bson.ObjectId) *KeyExistsMiddleware {
	return &KeyExistsMiddleware{manager, oAuthServerID}
}

// Handler is the middleware method.
func (m *KeyExistsMiddleware) Handler(handler http.Handler) http.Handler {
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
			panic(ErrAuthorizationFieldNotFound)
		}

		if strings.ToLower(parts[0]) != "bearer" {
			logger.Info("Bearer token malformed")
			panic(ErrBearerMalformed)
		}

		accessToken := parts[1]
		thisSessionState, keyExists := m.CheckSessionAndIdentityForValidKey(accessToken)

		if !keyExists {
			log.WithFields(log.Fields{
				"path":   r.RequestURI,
				"origin": r.RemoteAddr,
				"key":    accessToken,
			}).Info("Attempted access with non-existent key.")
			panic(ErrAccessTokenNotAuthorized)
		}

		if m.oAuthServerID != thisSessionState.OAuthServerID {
			log.WithFields(log.Fields{
				"path":   r.RequestURI,
				"origin": r.RemoteAddr,
				"key":    accessToken,
				"sessionOAuthServerID": thisSessionState.OAuthServerID,
				"authOAuthServerID":    m.oAuthServerID,
			}).Info("Attempted access with the key issued by other OAuth provider.")
			panic(ErrAccessTokenOfOtherOrigin)
		}

		ctx := context.WithValue(r.Context(), SessionData, thisSessionState)
		ctx = context.WithValue(ctx, AuthHeaderValue, accessToken)

		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}

// CheckSessionAndIdentityForValidKey ensures we have the valid key in the session store
func (m *KeyExistsMiddleware) CheckSessionAndIdentityForValidKey(key string) (session.SessionState, bool) {
	var thisSession session.SessionState

	// Checks if the key is present on the cache and if it didn't expire yet
	log.Debug("Querying keystore")
	exists, err := m.manager.KeyExists(key)
	if nil != err {
		panic(err)
	}

	if !exists {
		log.Debug("Key not found in keystore")
		return thisSession, false
	}

	// 2. If not there, get it from the AuthorizationHandler
	return m.manager.IsKeyAuthorised(key)
}
