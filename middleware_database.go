package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

// Database represents a database connection
type Database struct {
	dba *DatabaseAccessor
}

//Middleware is the gin middleware for the database, this is different from the others since it's a
//database handler
func (m *Database) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Debug("Starting Database middleware")

		reqSession := m.dba.Clone()
		defer reqSession.Close()
		m.dba.Set(c, reqSession)
	}
}
