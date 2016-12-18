package auth

import (
	"errors"
	"net/http"

	"context"

	jwt "github.com/dgrijalva/jwt-go"
)

type JWTMiddleware struct {
	Config JWTConfig
}

func NewJWTMiddleware(config JWTConfig) *JWTMiddleware {
	return &JWTMiddleware{config}
}

// Handler implementation
func (m *JWTMiddleware) Handler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		parser := JWTParser{m.Config}
		token, err := parser.Parse(r)

		if err != nil {
			m.Config.Unauthorized(w, r, err)
			return
		}

		claims := token.Claims.(jwt.MapClaims)

		id := claims["id"].(string)
		context.WithValue(r.Context(), "JWT_PAYLOAD", claims)
		context.WithValue(r.Context(), "userID", id)

		if !m.Config.Authorizator(id, w, r) {
			m.Config.Unauthorized(w, r, errors.New("you don't have permission to access"))
			return
		}

		handler.ServeHTTP(w, r)
	})
}
