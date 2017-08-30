package jwt

import (
	"net/http"
	"time"

	"github.com/hellofresh/janus/pkg/config"
)

// Guard struct
type Guard struct {
	ParserConfig

	// User can define own Unauthorized func.
	Unauthorized func(w http.ResponseWriter, r *http.Request, err error)

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
	}
}
