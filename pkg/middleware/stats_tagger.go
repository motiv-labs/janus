package middleware

import (
	"net/http"

	"go.opencensus.io/tag"
)

// statsTagger is a middleware that takes a list of tags and adds them into context to be propagated
type statsTagger struct {
	tags []tag.Mutator
}

// NewStatsTagger creates a new instance of statsTagger
func NewStatsTagger(tags []tag.Mutator) *statsTagger {
	metricKeyInserter := &statsTagger{tags}
	return metricKeyInserter
}

// Handler is the middleware function
func (h *statsTagger) Handler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		for _, t := range h.tags {
			ctx, _ = tag.New(ctx, t)
		}

		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}
