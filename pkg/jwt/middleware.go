package jwt

import (
	"context"
	"errors"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
)

// Payload Represents the context key
type Payload struct{}

// UserID Represents the user context key
type UserID struct{}

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
		token, err := parser.ParseFromRequest(r)

		if err != nil {
			m.Guard.Unauthorized(w, r, err)
			return
		}

		claims := token.Claims.(jwt.MapClaims)

		id := claims["id"].(string)
		context.WithValue(r.Context(), Payload{}, claims)
		context.WithValue(r.Context(), UserID{}, id)

		if !m.Guard.Authorizator(id, w, r) {
			m.Guard.Unauthorized(w, r, errors.New("you don't have permission to access"))
			return
		}

		handler.ServeHTTP(w, r)
	})
}
