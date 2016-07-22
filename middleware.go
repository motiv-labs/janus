package main

import (
	"github.com/kataras/iris"
	"github.com/valyala/fasthttp"
	log "github.com/Sirupsen/logrus"
)

// Middleware wraps up the APIDefinition object to be included in a
// middleware handler, this can probably be handled better.
type Middleware struct {
	Spec *APISpec
}

type MiddlewareImplementation interface {
	ProcessRequest(req fasthttp.Request, resp fasthttp.Response, c *iris.Context) (error, int)
}

// Generic middleware caller to make extension easier
func CreateMiddleware(mw MiddlewareImplementation, tykMwSuper *Middleware) {
	irisHandler := func(c *iris.Context) {
		req := c.Request
		res := c.Response

		reqErr, errCode := mw.ProcessRequest(req, res, c)

		if reqErr != nil {
			c.JSON(errCode, reqErr.Error())
			return
		}

		c.SetStatusCode(errCode)
		c.Next()
	}

	iris.UseFunc(irisHandler)
}

func (o Middleware) CheckSessionAndIdentityForValidKey(key string) (SessionState, bool) {
	var thisSession SessionState
	oAuthManager := o.Spec.OAuthManager

	//Checks if the key is present on the cache and if it didn't expire yet
	log.Debug("Querying keystore")
	if !oAuthManager.KeyExists(key) {
		log.Debug("Key not found in keystore")
		return thisSession, false
	}

	// 2. If not there, get it from the AuthorizationHandler
	thisSession, found := o.Spec.OAuthManager.IsKeyAuthorised(key)
	if found {
		// If not in Session, and got it from AuthHandler, create a session with a new TTL
		log.Info("Recreating session for key: ", key)

		// cache it
		o.Spec.OAuthManager.Set(key, thisSession, thisSession.ExpiresIn)
	}

	return thisSession, found
}
