package oauth2

import (
	"context"
	"errors"

	jwtBase "github.com/golang-jwt/jwt/v5"
	"github.com/hellofresh/janus/pkg/jwt"
	"github.com/hellofresh/janus/pkg/metrics"
	obs "github.com/hellofresh/janus/pkg/observability"
	"github.com/hellofresh/stats-go/bucket"
	"github.com/hellofresh/stats-go/client"
	log "github.com/sirupsen/logrus"
	"go.opencensus.io/stats"
	"go.opencensus.io/tag"
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

		switch {
		case errors.Is(err, jwtBase.ErrTokenExpired):
			shouldReport(ctx, stats, "ValidationErrorExpired")
			return false
		case errors.Is(err, jwtBase.ErrTokenInvalidClaims):
			shouldReport(ctx, stats, "ValidationErrorClaimsInvalid")
			return false
		case errors.Is(err, jwtBase.ErrTokenUsedBeforeIssued):
			shouldReport(ctx, stats, "ValidationErrorIssuedAt")
			return false
		case errors.Is(err, jwtBase.ErrTokenNotValidYet):
			shouldReport(ctx, stats, "ValidationErrorNotValidYet")
			return false
		case errors.Is(err, jwtBase.ErrTokenInvalidIssuer):
			shouldReport(ctx, stats, "ValidationErrorIssuer")
			return false
		case errors.Is(err, jwtBase.ErrTokenMalformed):
			shouldReport(ctx, stats, "ValidationErrorMalformed")
			return false
		case errors.Is(err, jwtBase.ErrSignatureInvalid):
			shouldReport(ctx, stats, "ValidationErrorSignatureInvalid")
			return false
		case errors.Is(err, jwtBase.ErrTokenUnverifiable):
			shouldReport(ctx, stats, "ValidationErrorUnverifiable")
			return false
		default:
			shouldReport(ctx, stats, "ErrFailedToParse")
			return false
		}
	}

	return true
}

func shouldReport(ctx context.Context, client client.Client, operation string) {
	client.TrackMetric("tokens", bucket.MetricOperation{"jwt-manager", "parse-error", operation})

	// OpenCensus stats
	ctx, _ = tag.New(ctx, tag.Insert(obs.KeyJWTValidationErrorType, operation))
	stats.Record(ctx, obs.MJWTManagerValidationErrors.M(1))
}
