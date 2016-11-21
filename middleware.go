package janus

import (
	"net/http"

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
