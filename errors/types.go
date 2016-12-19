package errors

import "net/http"

var (
	ErrClientIdNotFound = New(http.StatusBadRequest, "Invalid client ID provided")
	// ErrInvalidID represents an invalid ID
	ErrInvalidID = New(http.StatusBadRequest, "Please provide a valid ID")
)
