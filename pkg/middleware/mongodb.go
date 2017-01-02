package middleware

import (
	"context"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/hellofresh/janus/pkg/request"
)

var (
	// ContextKeyDatabase defines the db context key
	ContextKeyDatabase = request.ContextKey("db")
)

type MongoDB struct {
	accessor *DatabaseAccessor
}

func NewMongoDB(accessor *DatabaseAccessor) *MongoDB {
	return &MongoDB{accessor}
}

func (m *MongoDB) Handler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Starting Database middleware")

		reqSession := m.accessor.Clone()
		defer reqSession.Close()

		ctx := context.WithValue(r.Context(), ContextKeyDatabase, m.accessor.DB(""))
		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}
