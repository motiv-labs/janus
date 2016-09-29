package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Middleware wraps up the APIDefinition object to be included in a
// middleware handler, this can probably be handled better.
type Middleware struct {
	Spec   *APISpec
	Logger *Logger
}

type MiddlewareImplementation interface {
	ProcessRequest(req *http.Request, c *gin.Context) (error, int)
}

// Generic middleware caller to make extension easier
func CreateMiddleware(mw MiddlewareImplementation) gin.HandlerFunc {
	return func(c *gin.Context) {
		reqErr, errCode := mw.ProcessRequest(c.Request, c)

		if reqErr != nil {
			c.Abort()
			c.JSON(errCode, reqErr.Error())
		} else {
			c.Next()
			c.String(errCode, "")
		}
	}
}

func (o Middleware) CheckSessionAndIdentityForValidKey(key string) (SessionState, bool) {
	var thisSession SessionState
	oAuthManager := o.Spec.OAuthManager

	//Checks if the key is present on the cache and if it didn't expire yet
	o.Logger.Debug("Querying keystore")
	if !oAuthManager.KeyExists(key) {
		o.Logger.Debug("Key not found in keystore")
		return thisSession, false
	}

	// 2. If not there, get it from the AuthorizationHandler
	return o.Spec.OAuthManager.IsKeyAuthorised(key)
}
