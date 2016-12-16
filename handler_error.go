package janus

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/hellofresh/janus/errors"
)

// ErrInvalidID represents an invalid ID
var ErrInvalidID = errors.New(http.StatusBadRequest, "Please provide a valid ID")

// RecoveryHandler handler for the apis
func RecoveryHandler(rw http.ResponseWriter, err interface{}) {
	switch err.(type) {
	case *errors.Error:
		internalErr := err.(*errors.Error)
		log.Error(internalErr.Error())
		// c.JSON(internalErr.Code, internalErr.Error())
	default:
		// c.JSON(http.StatusInternalServerError, err)
	}
}
