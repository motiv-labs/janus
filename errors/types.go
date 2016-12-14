package errors

import "net/http"

var (
	ErrClientIdNotFound = New(http.StatusBadRequest, "Invalid client ID provided")
)
