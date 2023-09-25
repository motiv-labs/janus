package basic

import (
	"net/http"

	"github.com/hellofresh/janus/pkg/errors"
	"github.com/hellofresh/janus/pkg/plugin/basic/encrypt"
	log "github.com/sirupsen/logrus"
)

// NewBasicAuth is a HTTP basic auth middleware
func NewBasicAuth(repo Repository) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Debug("Starting basic auth middleware")
			logger := log.WithFields(log.Fields{
				"path":   r.RequestURI,
				"origin": r.RemoteAddr,
			})

			username, password, authOK := r.BasicAuth()
			if !authOK {
				errors.Handler(w, r, ErrNotAuthorized)
				return
			}

			var found bool
			users, err := repo.FindAll()
			if err != nil {
				log.WithError(err).Error("Error when getting all users")
				errors.Handler(w, r, errors.New(http.StatusInternalServerError, "there was an error when looking for users"))
				return
			}

			hash := encrypt.Hash{}

			for _, u := range users {
				//if username == u.Username && (subtle.ConstantTimeCompare([]byte(password), []byte(u.Password)) == 1) {
				if username == u.Username && (hash.Compare(u.Password, password) == nil) {
					found = true
					break
				}
			}

			if !found {
				logger.Debug("Invalid user/password provided.")
				errors.Handler(w, r, ErrNotAuthorized)
				return
			}

			handler.ServeHTTP(w, r)
		})
	}
}
