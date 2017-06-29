package oauth2

import (
	"context"
	"net/http"
	"strings"

	"github.com/hellofresh/janus/pkg/errors"
	"github.com/hellofresh/janus/pkg/oauth"
	"github.com/hellofresh/janus/pkg/request"
	log "github.com/sirupsen/logrus"
)

// Enums for keys to be stored in a session context - this is how gorilla expects
// these to be implemented and is lifted pretty much from docs
var (
	SessionData     = request.ContextKey("session_data")
	AuthHeaderValue = request.ContextKey("auth_header")
)

// NewKeyExistsMiddleware creates a new instance of KeyExistsMiddleware
func NewKeyExistsMiddleware(manager oauth.Manager) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
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
				panic(errors.ErrAuthorizationFieldNotFound)
			}

			if strings.ToLower(parts[0]) != "bearer" {
				logger.Warn("Bearer token malformed")
				panic(errors.ErrBearerMalformed)
			}

			accessToken := parts[1]
			thisSessionState, keyExists := manager.IsKeyAuthorised(accessToken)

			if !keyExists {
				log.WithFields(log.Fields{
					"path":   r.RequestURI,
					"origin": r.RemoteAddr,
					"key":    accessToken,
				}).Debug("Attempted access with invalid key.")
				panic(errors.ErrAccessTokenNotAuthorized)
			}

			ctx := context.WithValue(r.Context(), SessionData, thisSessionState)
			ctx = context.WithValue(ctx, AuthHeaderValue, accessToken)

			handler.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
