package company

import (
	"github.com/hellofresh/janus/pkg/errors"
	log "github.com/sirupsen/logrus"
	"net/http"
)



const companyHeader = "X-Company"
// NewCompany is a HTTP company middleware
func NewCompany(company Company, repo Repository) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			query := r.URL.Query()
			log.Debugf("this is company %v", company)
			username, _, ok := r.BasicAuth(); if !ok {
				errors.Handler(w, r, ErrNotAuthorized)
			}
			log.Debugf("got username %s", username)
			c, err := repo.FindAll()
			if err != nil {
				log.Errorf("find company by username failed: %v", err)
			}
			log.Debugf("got company %s", username)

			if c != nil {
				for _, u := range c {
					if u.Username == username {
						company.Company = u.Company
						break
					}
				}

				// if the header already exists, delete it and write a new one it
				if company.Company != "" {
					if r.Header.Get(companyHeader) != "" {
						r.Header.Del(companyHeader)
					}
					r.Header.Add(companyHeader, company.Company)
				} else {
					log.Debugf("No company associated with user")
				}
			}

			r.URL.RawQuery = query.Encode()
			next.ServeHTTP(w, r)
		})
	}
}
