package oauth

import "net/http"
import "github.com/hellofresh/janus/pkg/errors"

var (
	// ErrClientIDNotFound is raised when a client_id was not found
	ErrClientIDNotFound = errors.New(http.StatusBadRequest, "Invalid client ID provided")

	// ErrAccessTokenOfOtherOrigin is used when the access token is of other origin
	ErrAccessTokenOfOtherOrigin = errors.New(http.StatusUnauthorized, "access token of other origin")

	// ErrOauthServerNotFound is used when the oauth server was not found in the datastore
	ErrOauthServerNotFound = errors.New(http.StatusNotFound, "oauth server not found")

	// ErrDBContextNotSet is used when the database request context is not set
	ErrDBContextNotSet = errors.New(http.StatusInternalServerError, "DB context was not set for this request")

	// ErrJWTSecretMissing is used when the database request context is not set
	ErrJWTSecretMissing = errors.New(http.StatusBadRequest, "You need to set a JWT secret")

	// ErrUnknownManager is used when a manager type is not known
	ErrUnknownManager = errors.New(http.StatusBadRequest, "Unknown manager type provided")

	// ErrUnknownStrategy is used when a token strategy is not known
	ErrUnknownStrategy = errors.New(http.StatusBadRequest, "Unknown token strategy type provided")
)
