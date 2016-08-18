package main

import (
	"github.com/hellofresh/api-gateway/storage"
	"github.com/gin-gonic/gin"
	log "github.com/Sirupsen/logrus"
)

// a silly example
type Database struct {
	dba *storage.DatabaseAccessor
}

//Middleware is the gin middleware for the database, this is different from the others since it's a
//database handler
func (m Database) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Debug("Starting Database middleware")

		reqSession := m.dba.Clone()
		defer reqSession.Close()
		m.dba.Set(c, reqSession)
	}
}
