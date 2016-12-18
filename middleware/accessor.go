package middleware

import (
	log "github.com/Sirupsen/logrus"
	"gopkg.in/mgo.v2"
)

// DatabaseAccessor represents a mongo database encapsulation
type DatabaseAccessor struct {
	*mgo.Session
}

// InitDB creates a new mongo db server
func InitDB(dsn string) (*DatabaseAccessor, error) {
	log.Debugf("Trying to connect to %s", dsn)
	session, err := mgo.Dial(dsn)

	if err == nil {
		log.Debug("Connected to mongodb")
		session.SetMode(mgo.Monotonic, true)
		return &DatabaseAccessor{session}, nil
	}

	return nil, err
}
