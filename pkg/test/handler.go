package test

import (
	"net/http"

	"github.com/hellofresh/janus/pkg/errors"
	"github.com/hellofresh/janus/pkg/response"
)

// Ping is a test handler
func Ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte("OK\n"))
}

// RecoveryHandler represents the recovery handler
func RecoveryHandler(w http.ResponseWriter, r *http.Request, err interface{}) {
	switch internalErr := err.(type) {
	case *errors.Error:
		response.JSON(w, internalErr.Code, internalErr.Error())
	default:
		response.JSON(w, http.StatusInternalServerError, err)
	}
}
