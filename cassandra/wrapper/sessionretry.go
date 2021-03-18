package wrapper

import (
	"github.com/gocql/gocql"
	log "github.com/sirupsen/logrus"
)

// sessionRetry is an implementation of SessionInterface
type sessionRetry struct {
	goCqlSession *gocql.Session
}

// Query wrapper to be able to return our own QueryInterface
func (s sessionRetry) Query( stmt string, values ...interface{}) QueryInterface {
	log.Debug("running SessionRetry Query() method")

	return queryRetry{goCqlQuery: s.goCqlSession.Query(stmt, values...)}
}

// Close wrapper to be able to run goCql method
func (s sessionRetry) Close() {
	log.Debug("running SessionRetry Close() method")

	s.goCqlSession.Close()
}
