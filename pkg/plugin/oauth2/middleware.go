package oauth2

import (
	"context"
	"net/http"
	"strings"

	"github.com/hellofresh/janus/pkg/errors"
	"github.com/hellofresh/janus/pkg/metrics"
	obs "github.com/hellofresh/janus/pkg/observability"
	"github.com/hellofresh/stats-go/bucket"
	log "github.com/sirupsen/logrus"
	"go.opencensus.io/stats"
)

const (
	tokensSection = "tokens"
)

// Enums for keys to be stored in a session context - this is how gorilla expects
// these to be implemented and is lifted pretty much from docs
var (
	AuthHeaderValue = ContextKey("auth_header")

	// ErrAuthorizationFieldNotFound is used when the http Authorization header is missing from the request
	ErrAuthorizationFieldNotFound = errors.New(http.StatusBadRequest, "authorization field missing")
	// ErrBearerMalformed is used when the Bearer string in the Authorization header is not found or is malformed
	ErrBearerMalformed = errors.New(http.StatusBadRequest, "bearer token malformed")
	// ErrAccessTokenNotAuthorized is used when the access token is not found on the storage
	ErrAccessTokenNotAuthorized = errors.New(http.StatusUnauthorized, "access token not authorized")
)

// ContextKey is used to create context keys that are concurrent safe
type ContextKey string

func (c ContextKey) String() string {
	return "janus." + string(c)
}

// NewKeyExistsMiddleware creates a new instance of KeyExistsMiddleware
func NewKeyExistsMiddleware(manager Manager) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Debug("Starting Oauth2KeyExists middleware")
			statsClient := metrics.WithContext(r.Context())

			logger := log.WithFields(log.Fields{
				"path":   r.RequestURI,
				"origin": r.RemoteAddr,
			})

			// We're using OAuth, start checking for access keys
			authHeaderValue := r.Header.Get("Authorization")
			parts := strings.Split(authHeaderValue, " ")
			if len(parts) < 2 {
				logger.Warn("Attempted access with malformed header, no auth header found.")
				statsClient.TrackOperation(tokensSection, bucket.MetricOperation{"key-exists", "header"}, nil, false)
				stats.Record(r.Context(), obs.MOAuth2MissingHeader.M(1))
				errors.Handler(w, ErrAuthorizationFieldNotFound)
				return
			}
			statsClient.TrackOperation(tokensSection, bucket.MetricOperation{"key-exists", "header"}, nil, true)

			if strings.ToLower(parts[0]) != "bearer" {
				logger.Warn("Bearer token malformed")
				statsClient.TrackOperation(tokensSection, bucket.MetricOperation{"key-exists", "malformed"}, nil, false)
				stats.Record(r.Context(), obs.MOAuth2MalformedHeader.M(1))
				errors.Handler(w, ErrBearerMalformed)
				return
			}
			statsClient.TrackOperation(tokensSection, bucket.MetricOperation{"key-exists", "malformed"}, nil, true)

			accessToken := parts[1]
			keyExists := manager.IsKeyAuthorized(r.Context(), accessToken)
			statsClient.TrackOperation(tokensSection, bucket.MetricOperation{"key-exists", "authorized"}, nil, keyExists)
			if keyExists {
				stats.Record(r.Context(), obs.MOAuth2Authorized.M(1))
			} else {
				stats.Record(r.Context(), obs.MOAuth2Unauthorized.M(1))
			}

			if !keyExists {
				log.WithFields(log.Fields{
					"path":   r.RequestURI,
					"origin": r.RemoteAddr,
					"key":    accessToken,
				}).Debug("Attempted access with invalid key.")
				errors.Handler(w, ErrAccessTokenNotAuthorized)
				return
			}

			ctx := context.WithValue(r.Context(), AuthHeaderValue, accessToken)
			handler.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
