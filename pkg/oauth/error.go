package oauth

import "net/http"
import "github.com/hellofresh/janus/pkg/errors"

var (
	// ErrClientIDNotFound is raised when a client_id was not found
	ErrClientIDNotFound = errors.New(http.StatusBadRequest, "Invalid client ID provided")

	// ErrAuthorizationFieldNotFound is used when the http Authorization header is missing from the request
	ErrAuthorizationFieldNotFound = errors.New(http.StatusBadRequest, "authorization field missing")

	// ErrBearerMalformed is used when the Bearer string in the Authorization header is not found or is malformed
	ErrBearerMalformed = errors.New(http.StatusBadRequest, "bearer token malformed")

	// ErrAccessTokenNotAuthorized is used when the access token is not found on the storage
	ErrAccessTokenNotAuthorized = errors.New(http.StatusUnauthorized, "access token not authorized")

	// ErrOauthServerNotFound is used when the oauth server was not found in the datastore
	ErrOauthServerNotFound = errors.New(http.StatusNotFound, "oauth server not found")

	// ErrDBContextNotSet is used when the database request context is not set
	ErrDBContextNotSet = errors.New(http.StatusInternalServerError, "DB context was not set for this request")
)
