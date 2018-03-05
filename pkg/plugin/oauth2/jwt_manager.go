package oauth2

import (
	"context"

	jwtBase "github.com/dgrijalva/jwt-go"
	"github.com/hellofresh/janus/pkg/jwt"
	"github.com/hellofresh/janus/pkg/metrics"
	"github.com/hellofresh/stats-go/bucket"
	"github.com/hellofresh/stats-go/client"
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
func (m *JWTManager) IsKeyAuthorized(ctx context.Context, accessToken string) bool {
	if ctx == nil {
		return false
	}

	stats := metrics.WithContext(ctx)
	if stats == nil {
		return false
	}

	if _, err := m.parser.Parse(accessToken); err != nil {
		log.WithError(err).Info("Failed to parse and validate the JWT")

		switch jwtErr := err.(type) {
		case *jwtBase.ValidationError:
			shouldReport(stats, jwtErr.Errors&jwtBase.ValidationErrorExpired != 0, "ValidationErrorExpired")
			shouldReport(stats, jwtErr.Errors&jwtBase.ValidationErrorClaimsInvalid != 0, "ValidationErrorClaimsInvalid")
			shouldReport(stats, jwtErr.Errors&jwtBase.ValidationErrorIssuedAt != 0, "ValidationErrorIssuedAt")
			shouldReport(stats, jwtErr.Errors&jwtBase.ValidationErrorNotValidYet != 0, "ValidationErrorNotValidYet")
			shouldReport(stats, jwtErr.Errors&jwtBase.ValidationErrorIssuer != 0, "ValidationErrorIssuer")
			shouldReport(stats, jwtErr.Errors&jwtBase.ValidationErrorMalformed != 0, "ValidationErrorMalformed")
			shouldReport(stats, jwtErr.Errors&jwtBase.ValidationErrorSignatureInvalid != 0, "ValidationErrorSignatureInvalid")
			shouldReport(stats, jwtErr.Errors&jwtBase.ValidationErrorUnverifiable != 0, "ValidationErrorUnverifiable")
			return false
		default:
			stats.TrackMetric("tokens", bucket.MetricOperation{"jwt-manager", "parse-error", "ErrFailedToParse"})
			return false
		}
	}

	return true
}

func shouldReport(client client.Client, typeCheck bool, operation string) {
	if typeCheck {
		client.TrackMetric("tokens", bucket.MetricOperation{"jwt-manager", "parse-error", operation})
	}
}
