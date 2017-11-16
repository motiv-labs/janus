package oauth2

import (
	"net/http"

	"github.com/hellofresh/janus/pkg/errors"
)

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

	// ErrInvalidIntrospectionURL is used when an introspection URL is invalid
	ErrInvalidIntrospectionURL = errors.New(http.StatusBadRequest, "The provided introspection URL is invalid")

	// ErrOauthServerNameExists is used when the Oauth Server name is already registered on the datastore
	ErrOauthServerNameExists = errors.New(http.StatusConflict, "oauth server name is already registered")
)
