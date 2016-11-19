package main

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

// Middleware wraps up the APIDefinition object to be included in a
// middleware handler, this can probably be handled better.
type Middleware struct {
	Spec *APISpec
}

type MiddlewareImplementation interface {
	ProcessRequest(req *http.Request, c *gin.Context) (error, int)
}

// Generic middleware caller to make extension easier
func CreateMiddleware(mw MiddlewareImplementation) gin.HandlerFunc {
	return func(c *gin.Context) {
		err, errCode := mw.ProcessRequest(c.Request, c)

		if err != nil {
			c.Abort()
			c.JSON(errCode, err.Error())
		} else {
			c.Next()
		}
	}
}

func (o *Middleware) CheckSessionAndIdentityForValidKey(key string) (SessionState, bool) {
	var thisSession SessionState
	oAuthManager := o.Spec.OAuthManager

	//Checks if the key is present on the cache and if it didn't expire yet
	log.Debug("Querying keystore")
	if !oAuthManager.KeyExists(key) {
		log.Debug("Key not found in keystore")
		return thisSession, false
	}

	// 2. If not there, get it from the AuthorizationHandler
	return o.Spec.OAuthManager.IsKeyAuthorised(key)
}
