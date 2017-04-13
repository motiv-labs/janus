package api

import "net/http"
import "github.com/hellofresh/janus/pkg/errors"

var (
	// ErrAPIDefinitionNotFound is used when the api was not found in the datastore
	ErrAPIDefinitionNotFound = errors.New(http.StatusNotFound, "api definition not found")

	// ErrAPINameExists is used when the API name is already registeres on the datastore
	ErrAPINameExists = errors.New(http.StatusBadRequest, "api name is already registered")

	// ErrAPIListenPathExists is used when the API name is already registeres on the datastore
	ErrAPIListenPathExists = errors.New(http.StatusBadRequest, "api listen path is already registered")

	// ErrDBContextNotSet is used when the database request context is not set
	ErrDBContextNotSet = errors.New(http.StatusInternalServerError, "DB context was not set for this request")
)
