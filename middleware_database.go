package main

import (
	"github.com/hellofresh/api-gateway/storage"
	"net/http"
	"github.com/gin-gonic/gin"
)

// a silly example
type Database struct {
	*Middleware
	dba *storage.DatabaseAccessor
}

//Important staff, iris middleware must implement the iris.Handler interface which is:
func (m Database) ProcessRequest(req *http.Request, c *gin.Context) (error, int) {
	m.Logger.Debug("Starting Database middleware")

	reqSession := m.dba.Clone()
	defer reqSession.Close()
	m.dba.Set(c, reqSession)

	return nil, http.StatusOK
}
