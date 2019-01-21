package bodylmt

import (
	"net/http"

	"code.cloudfoundry.org/bytefmt"
	"github.com/hellofresh/janus/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var (
	// ErrRequestEntityTooLarge is thrown when a body size is bigger then the limit specified
	ErrRequestEntityTooLarge = errors.New(http.StatusRequestEntityTooLarge, http.StatusText(http.StatusRequestEntityTooLarge))
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
				errors.Handler(w, r, ErrRequestEntityTooLarge)
				return
			}

			r.Body = http.MaxBytesReader(w, r.Body, int64(limit))
			handler.ServeHTTP(w, r)
		})
	}
}
