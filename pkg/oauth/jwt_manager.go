package oauth

import (
	"time"

	"github.com/hellofresh/janus/pkg/jwt"
	log "github.com/sirupsen/logrus"
)

// JWTManager is responsible for managing the JWT tokens
type JWTManager struct {
	parser *jwt.Parser
}

// NewJWTManager creates a new instance of JWTManager
func NewJWTManager(parser *jwt.Parser) *JWTManager {
	return &JWTManager{parser}
}

// IsKeyAuthorized checks if the access token is valid
func (m *JWTManager) IsKeyAuthorized(accessToken string) bool {
	token, err := m.parser.Parse(accessToken)
	if err != nil {
		log.WithError(err).Info("Could not parse the JWT")
		return false
	}

	if claims, ok := m.parser.GetStandardClaims(token); ok && token.Valid {
		expiresAt := time.Unix(claims.ExpiresAt, 0)
		if time.Now().After(expiresAt) {
			return false
		}
	}

	return true
}
