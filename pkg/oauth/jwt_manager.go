package oauth

import (
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
	_, err := m.parser.Parse(accessToken)
	if err != nil {
		log.WithError(err).Info("Failed to parse and validate the JWT")
		return false
	}

	return true
}
