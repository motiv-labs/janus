package oauth

import (
	"time"

	"github.com/hellofresh/janus/pkg/jwt"
	"github.com/hellofresh/janus/pkg/session"
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

// Set returns nil since when we work with JWT we don't need to store them
func (m *JWTManager) Set(accessToken string, session session.State, resetTTLTo int64) error {
	return nil
}

// Remove returns nil because there is not storage to remove from
func (m *JWTManager) Remove(accessToken string) error {
	return nil
}

// IsKeyAuthorised checks if the access token is valid
func (m *JWTManager) IsKeyAuthorised(accessToken string) (session.State, bool) {
	var sessionState session.State

	token, err := m.parser.Parse(accessToken)
	if err != nil {
		log.WithError(err).Info("Failed to parse and validate the JWT")
		return sessionState, false
	}

	// as parser.Parse() does validation we are sure that token is valid at this point
	if claims, ok := m.parser.GetMapClaims(token); ok {
		sessionState.AccessToken = accessToken
		sessionState.ExpiresIn = claims["exp"].(int64)
	}

	return sessionState, true
}
