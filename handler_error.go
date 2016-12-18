package janus

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/hellofresh/janus/errors"
	"github.com/hellofresh/janus/response"
)

// RecoveryHandler handler for the apis
func RecoveryHandler(w http.ResponseWriter, r *http.Request, err interface{}) {
	switch err.(type) {
	case *errors.Error:
		internalErr := err.(*errors.Error)
		log.Error(internalErr.Error())
		response.JSON(w, internalErr.Code, internalErr.Error())
	default:
		response.JSON(w, http.StatusInternalServerError, err)
	}
}
