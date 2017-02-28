package store

import "net/http"
import "github.com/hellofresh/janus/pkg/errors"

var (
	// ErrUnknownStorage is used when the storage type is nknown
	ErrUnknownStorage = errors.New(http.StatusBadRequest, "Unknown storage type")
)
