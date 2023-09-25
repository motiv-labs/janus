/*
Package errors provides a nice way of handling http errors

Examples:
To create an error:

	err := errors.New(http.StatusBadRequest, "Something went wrong")
*/
package errors

import (
	"net/http"
	"runtime/debug"

	log "github.com/sirupsen/logrus"

	"github.com/hellofresh/janus/pkg/observability"
	"github.com/hellofresh/janus/pkg/render"
)

var (
	// ErrRouteNotFound happens when no route was matched
	ErrRouteNotFound = New(http.StatusNotFound, "no API found with those values")
	// ErrInvalidID represents an invalid identifier
	ErrInvalidID = New(http.StatusBadRequest, "please provide a valid ID")
)

// Error is a custom error that implements the `error` interface.
// When creating errors you should provide a code (could be and http status code)
// and a message, this way we can handle the errors in a centralized place.
type Error struct {
	Code    int    `json:"-"`
	Message string `json:"error"`
}

// New creates a new instance of Error
func New(code int, message string) *Error {
	return &Error{code, message}
}

func (e *Error) Error() string {
	return e.Message
}

// NotFound handler is called when no route is matched
func NotFound(w http.ResponseWriter, r *http.Request) {
	Handler(w, r, ErrRouteNotFound)
}

// RecoveryHandler handler is used when a panic happens
func RecoveryHandler(w http.ResponseWriter, r *http.Request, err interface{}) {
	Handler(w, r, err)
}

// Handler marshals an error to JSON, automatically escaping HTML and setting the
// Content-Type as application/json.
func Handler(w http.ResponseWriter, r *http.Request, err interface{}) {
	entry := log.WithField("request-id", observability.RequestIDFromContext(r.Context()))

	switch internalErr := err.(type) {
	case *Error:
		entry.WithFields(log.Fields{
			"code":       internalErr.Code,
			log.ErrorKey: internalErr.Error(),
		}).Info("Internal error handled")
		render.JSON(w, internalErr.Code, internalErr)
	case error:
		entry.WithError(internalErr).WithField("stack", string(debug.Stack())).Error("Internal server error handled")
		render.JSON(w, http.StatusInternalServerError, internalErr.Error())
	default:
		entry.WithFields(log.Fields{
			log.ErrorKey: err,
			"stack":      string(debug.Stack()),
		}).Error("Internal server error handled")
		render.JSON(w, http.StatusInternalServerError, err)
	}
}
