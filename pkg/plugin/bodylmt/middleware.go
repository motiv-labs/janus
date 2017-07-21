package bodylmt

import (
	"net/http"

	"code.cloudfoundry.org/bytefmt"

	log "github.com/sirupsen/logrus"
)

// NewBodyLimitMiddleware creates a new body limit middleware
func NewBodyLimitMiddleware(limit string) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.WithField("limit", limit).Debug("Starting body limit middleware")
			limit, err := bytefmt.ToBytes(limit)
			if err != nil {
				log.WithError(err).WithField("limit", limit).Error("invalid body-limit")
			}

			// Based on content length
			if r.ContentLength > int64(limit) {
				w.WriteHeader(http.StatusRequestEntityTooLarge)
				return
			}

			r.Body = http.MaxBytesReader(w, r.Body, int64(limit))
			handler.ServeHTTP(w, r)
		})
	}
}
