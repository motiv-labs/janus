package web

import (
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/hellofresh/janus/pkg/errors"
	"github.com/hellofresh/janus/pkg/response"
)

func Home() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response.JSON(w, http.StatusOK, fmt.Sprintf("Welcome to Janus"))
	}
}

// RecoveryHandler handler for the apis
func RecoveryHandler(w http.ResponseWriter, r *http.Request, err interface{}) {
	switch internalErr := err.(type) {
	case *errors.Error:
		log.Error(internalErr.Error())
		response.JSON(w, internalErr.Code, internalErr)
	case error:
		jsonErr := errors.New(http.StatusInternalServerError, internalErr.Error())
		response.JSON(w, jsonErr.Code, jsonErr)
	}
}

func Heartbeat() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response.JSON(w, http.StatusOK, nil)
	}
}
