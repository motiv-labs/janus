package organization

import (
	"encoding/json"
	"github.com/hellofresh/janus/pkg/errors"
	"github.com/hellofresh/janus/pkg/plugin/basic/encrypt"
	log "github.com/sirupsen/logrus"
	"net/http"
)

const organizationHeader = "X-Organization"
const organizationConfigHeader = "X-OrganizationConfig"

// NewOrganization is a HTTP organization middleware
func NewOrganization(organization Organization, repo Repository) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			query := r.URL.Query()

			log.Debug("Starting organization auth middleware")
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
			var organizationConfig *OrganizationConfig

			for _, u := range users {
				//if username == u.Username && (subtle.ConstantTimeCompare([]byte(password), []byte(u.Password)) == 1) {
				if username == u.Username && (hash.Compare(u.Password, password) == nil) {
					found = true
					organization.Organization = u.Organization
					organizationConfig, err = repo.FindOrganization(u.Organization)
					if err != nil {
						log.WithError(err).Error("Error finding organization")
						errors.Handler(w, r, errors.New(http.StatusInternalServerError, "there was an error when looking for organization"))
						return
					}
					break
				}
			}

			if !found {
				logger.Debug("Invalid user/password provided.")
				errors.Handler(w, r, ErrNotAuthorized)
				return
			}

			// if the header already exists, delete it and write a new one it
			if organization.Organization != "" {
				if r.Header.Get(organizationHeader) != "" {
					r.Header.Del(organizationHeader)
				}
				r.Header.Add(organizationHeader, organization.Organization)

				if r.Header.Get(organizationConfigHeader) != "" {
					r.Header.Del(organizationConfigHeader)
				}
				bOrganizationConfig, err := json.Marshal(organizationConfig)
				if err != nil {
					log.WithError(err).Error("Error marshaling organization for config header")
					errors.Handler(w, r, errors.New(http.StatusInternalServerError, "there was an error when setting config header"))
					return
				}
				r.Header.Add(organizationConfigHeader, string(bOrganizationConfig))
			} else {
				log.Debugf("No organization associated with user")
			}

			r.URL.RawQuery = query.Encode()
			next.ServeHTTP(w, r)
		})
	}
}
