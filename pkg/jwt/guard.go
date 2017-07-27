package jwt

import (
	"net/http"
	"time"

	"github.com/hellofresh/janus/pkg/config"
	"github.com/hellofresh/janus/pkg/errors"
)

// Guard struct
type Guard struct {
	ParserConfig

	// Callback function that should perform the authentication of the user based on userID and
	// password. Must return true on success, false on failure. Required.
	// Option return user id, if so, user id will be stored in Claim Array.
	Authenticator func(userID string, password string) (string, bool)

	// User can define own Unauthorized func.
	Unauthorized func(w http.ResponseWriter, r *http.Request, err error)

	// Callback function that should perform the authorization of the authenticated user. Called
	// only after an authentication success. Must return true on success, false on failure.
	// Optional, default to success.
	Authorizator func(userID string, w http.ResponseWriter, r *http.Request) bool

	// Duration that a jwt token is valid. Optional, defaults to one hour.
	Timeout time.Duration

	// SigningMethod defines new token signing algorithm/key pair.
	SigningMethod SigningMethod

	// This field allows clients to refresh their token until MaxRefresh has passed.
	// Note that clients can refresh their token in the last moment of MaxRefresh.
	// This means that the maximum validity timespan for a token is MaxRefresh + Timeout.
	// Optional, defaults to 0 meaning not refreshable.
	MaxRefresh time.Duration
}

// NewGuard creates a new instance of Guard with default handlers
func NewGuard(cred config.Credentials) Guard {
	return Guard{
		ParserConfig: ParserConfig{
			SigningMethods: []SigningMethod{{Alg: cred.Algorithm, Key: cred.Secret}},
			TokenLookup:    "header:Authorization",
		},
		SigningMethod: SigningMethod{Alg: cred.Algorithm, Key: cred.Secret},
		Timeout:       time.Hour,
		MaxRefresh:    time.Hour * 24,
		Authenticator: func(userID string, password string) (string, bool) {
			if userID == cred.Username && password == cred.Password {
				return userID, true
			}

			return userID, false
		},
		Authorizator: func(userID string, w http.ResponseWriter, r *http.Request) bool {
			return userID == cred.Username
		},
		Unauthorized: func(w http.ResponseWriter, r *http.Request, err error) {
			errors.Handler(w, errors.New(http.StatusUnauthorized, err.Error()))
		},
	}
}
