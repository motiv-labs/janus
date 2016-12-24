package api

import "net/http"
import "github.com/hellofresh/janus/errors"

var (
	// ErrAPIDefinitionNotFound is used when the api was not found in the datastore
	ErrAPIDefinitionNotFound = errors.New(http.StatusNotFound, "api definition not found")

	// ErrDBContextNotSet is used when the database request context is not set
	ErrDBContextNotSet = errors.New(http.StatusInternalServerError, "DB context was not set for this request")
)
