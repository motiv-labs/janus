package oauth2

import (
	"net/http"

	"github.com/hellofresh/janus/pkg/jwt"
	"github.com/hellofresh/janus/pkg/oauth"
	log "github.com/sirupsen/logrus"
)

// NewRevokeRulesMiddleware creates a new revoke rules middleware
func NewRevokeRulesMiddleware(parser *jwt.Parser, accessRules []*oauth.AccessRule) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Debug("Starting revoke rules middleware")

			// If no rules are set then lets not parse the token to avoid performance issues
			if len(accessRules) <= 0 {
				handler.ServeHTTP(w, r)
				return
			}

			token, err := parser.ParseFromRequest(r)
			if err != nil {
				log.WithError(err).Debug("Could not parse the JWT")
				handler.ServeHTTP(w, r)
				return
			}

			if claims, ok := parser.GetMapClaims(token); ok && token.Valid {
				for _, rule := range accessRules {
					allowed, err := rule.IsAllowed(claims)
					if err != nil {
						log.WithError(err).Debug("Rule is not allowed")
						continue
					}

					if allowed {
						handler.ServeHTTP(w, r)
					} else {
						w.WriteHeader(http.StatusUnauthorized)
						return
					}
				}
			}

			handler.ServeHTTP(w, r)
		})
	}
}
