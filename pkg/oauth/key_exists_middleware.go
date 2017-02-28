package oauth

import (
	"context"
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/hellofresh/janus/pkg/request"
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
	manager Manager
	// TODO pass api.Spec pointer here, currently impossible as circullar referens occurs
	oAuthServerID bson.ObjectId
}

// NewKeyExistsMiddleware creates a new instance of KeyExistsMiddleware
func NewKeyExistsMiddleware(manager Manager, oAuthServerID bson.ObjectId) *KeyExistsMiddleware {
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
			logger.Warn("Attempted access with malformed header, no auth header found.")
			panic(ErrAuthorizationFieldNotFound)
		}

		if strings.ToLower(parts[0]) != "bearer" {
			logger.Warn("Bearer token malformed")
			panic(ErrBearerMalformed)
		}

		accessToken := parts[1]
		thisSessionState, keyExists := m.manager.IsKeyAuthorised(accessToken)

		if !keyExists {
			log.WithFields(log.Fields{
				"path":   r.RequestURI,
				"origin": r.RemoteAddr,
				"key":    accessToken,
			}).Warn("Attempted access with non-existent key.")
			panic(ErrAccessTokenNotAuthorized)
		}

		ctx := context.WithValue(r.Context(), SessionData, thisSessionState)
		ctx = context.WithValue(ctx, AuthHeaderValue, accessToken)

		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}
