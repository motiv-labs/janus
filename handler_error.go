package main

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/hellofresh/janus/errors"
)

// ErrInvalidID represents an invalid ID
var ErrInvalidID = errors.New(http.StatusBadRequest, "Please provide a valid ID")

// Recovery handler for the apis
func recoveryHandler(c *gin.Context, err interface{}) {
	switch err.(type) {
	case *errors.Error:
		internalErr := err.(*errors.Error)
		log.Error(internalErr.Error())
		c.JSON(internalErr.Code, internalErr.Error())
	default:
		c.JSON(http.StatusInternalServerError, err)
	}
}
