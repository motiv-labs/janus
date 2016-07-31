package storage

import (
	log "github.com/Sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"github.com/gin-gonic/gin"
)

// DatabaseAccessor represents a mongo database encapsulation
type DatabaseAccessor struct {
	*mgo.Session
	config Database
}

// NewServer create a new mongo db server
func NewServer(config Database) (*DatabaseAccessor, error) {
	log.Debugf("Trying to connect to %s", config.DSN)
	session, err := mgo.Dial(config.DSN)

	if err == nil {
		log.Debug("Connected to session")
		session.SetMode(mgo.Monotonic, true)
		return &DatabaseAccessor{session, config}, nil
	}

	return &DatabaseAccessor{}, err
}

// Set a session to a context
func (da *DatabaseAccessor) Set(c *gin.Context, session *mgo.Session) {
	db := da.DB(da.config.Name)
	c.Set("db", db)
	c.Set("mgoSession", session)
}
