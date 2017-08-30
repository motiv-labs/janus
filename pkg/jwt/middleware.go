package jwt

import (
	"net/http"
)

// Payload Represents the context key
type Payload struct{}

// User represents a logged in user
type User struct {
	Username string
	Email    string
}

// Middleware struct contains data and logic required for middleware functionality
type Middleware struct {
	Guard Guard
}

// NewMiddleware builds and returns new JWT middleware instance
func NewMiddleware(config Guard) *Middleware {
	return &Middleware{config}
}

// Handler implementation
func (m *Middleware) Handler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		parser := Parser{m.Guard.ParserConfig}
		_, err := parser.ParseFromRequest(r)
		if err != nil {
			m.Guard.Unauthorized(w, r, err)
			return
		}

		handler.ServeHTTP(w, r)
	})
}
