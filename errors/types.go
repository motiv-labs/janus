package errors

import "net/http"

var (
	// ErrInvalidID represents an invalid indentifier
	ErrInvalidID = New(http.StatusBadRequest, "Please provide a valid ID")
)
