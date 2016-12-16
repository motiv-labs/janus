package middleware

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/urfave/negroni"
)

//MongoSession is the gin middleware for the database, this is different from the others since it's a
//database handler
func MongoSession(accessor *DatabaseAccessor) negroni.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		log.Debug("Starting Database middleware")

		reqSession := accessor.Clone()
		defer reqSession.Close()
		accessor.Set(r, reqSession)
		next(rw, r)
	}
}
