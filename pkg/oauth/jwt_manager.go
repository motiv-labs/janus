package oauth

import (
	"fmt"

	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/dgrijalva/jwt-go"
	"github.com/hellofresh/janus/pkg/session"
)

// JWTManager is responsible for managing the JWT tokens
type JWTManager struct {
	Secret string
}

// Set returns nil since when we work with JWT we don't need to store them
func (m *JWTManager) Set(accessToken string, session session.SessionState, resetTTLTo int64) error {
	return nil
}

// Remove returns nil becuase there is not storage to remove from
func (m *JWTManager) Remove(accessToken string) error {
	return nil
}

// IsKeyAuthorised checks if the access token is valid
func (m *JWTManager) IsKeyAuthorised(accessToken string) (session.SessionState, bool) {
	var session session.SessionState

	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Errorf("Unexpected signing method: %v", token.Header["alg"])
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(m.Secret), nil
	})

	if err != nil {
		log.WithError(err).Error("Could not parse the JWT")
		return session, false
	}

	if claims, ok := token.Claims.(jwt.StandardClaims); ok && token.Valid {
		expiresAt := time.Unix(claims.ExpiresAt, 0)
		if time.Now().After(expiresAt) {
			return session, false
		}
		session.AccessToken = accessToken
		session.ExpiresIn = claims.ExpiresAt
	}

	return session, true
}
