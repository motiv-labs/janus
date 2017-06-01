package errors

import "net/http"

var (
	// ErrRouteNotFound happens when no route was matched
	ErrRouteNotFound = New(http.StatusNotFound, "no API found with those values")
	// ErrInvalidID represents an invalid identifier
	ErrInvalidID = New(http.StatusBadRequest, "please provide a valid ID")
	// ErrAuthorizationFieldNotFound is used when the http Authorization header is missing from the request
	ErrAuthorizationFieldNotFound = New(http.StatusBadRequest, "authorization field missing")
	// ErrBearerMalformed is used when the Bearer string in the Authorization header is not found or is malformed
	ErrBearerMalformed = New(http.StatusBadRequest, "bearer token malformed")
	// ErrAccessTokenNotAuthorized is used when the access token is not found on the storage
	ErrAccessTokenNotAuthorized = New(http.StatusUnauthorized, "access token not authorized")
	// ErrInvalidScheme is used when the access token is not found on the storage
	ErrInvalidScheme = New(http.StatusBadRequest, "scheme is not supported")
	// ErrInvalidPolicy is used when an invalid policy was provided
	ErrInvalidPolicy = New(http.StatusBadRequest, "policy is not supported")
	// ErrInvalidStorage is used when an invalid storage was provided
	ErrInvalidStorage = New(http.StatusBadRequest, "the storage that you are using is not supported for this feature")
)
