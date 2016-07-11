package storage

import (
	log "github.com/Sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"github.com/kataras/iris"
)

type DatabaseAccessor struct {
	*mgo.Session
	config Database
}

func NewServer(config Database) (*DatabaseAccessor, error) {
	log.Infof("Trying to connect to %s", config.DSN)
	session, err := mgo.Dial(config.DSN)

	if err == nil {
		log.Info("Connected to session")
		session.SetMode(mgo.Monotonic, true)
		return &DatabaseAccessor{session, config}, nil
	} else {
		return &DatabaseAccessor{}, err
	}
}

func (da *DatabaseAccessor) Set(c *iris.Context, session *mgo.Session) {
	db := da.DB(da.config.Name)
	c.Set("db", db)
	c.Set("mgoSession", session)
}
