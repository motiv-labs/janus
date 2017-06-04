package jwt

import (
	"context"
	"errors"
	"net/http"

	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
)

// Payload Represents the context key
type Payload struct{}

// UserID Represents the user context key
type UserID struct{}

// Middleware struct contains data and logic required for middleware functionality
type Middleware struct {
	Config Config
}

// NewMiddleware builds and returns new JWT middleware instance
func NewMiddleware(config Config) *Middleware {
	return &Middleware{config}
}

// Handler implementation
func (m *Middleware) Handler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		parser := Parser{m.Config}
		token, err := parser.Parse(r)

		if err != nil {
			m.Config.Unauthorized(w, r, err)
			return
		}

		claims := token.Claims.(jwt.MapClaims)

		id, err := m.getUserId(claims)
		if err != nil {
			m.Config.Unauthorized(w, r, errors.New("The user identifier is not found in token"))
		}

		context.WithValue(r.Context(), Payload{}, claims)
		context.WithValue(r.Context(), UserID{}, id)

		if !m.Config.Authorizator(id, w, r) {
			m.Config.Unauthorized(w, r, errors.New("you don't have permission to access"))
			return
		}

		handler.ServeHTTP(w, r)
	})
}

func (m *Middleware) getUserId(claims jwt.MapClaims) (string, error) {
	var userKey = m.Config.UserIDKey
	if "" == userKey {
		userKey = "sub"
	}

	id, ok := claims[userKey].(string)
	if !ok {
		return "", fmt.Errorf("User key \"%s\" not found in token payload", userKey)
	}

	return id, nil
}
