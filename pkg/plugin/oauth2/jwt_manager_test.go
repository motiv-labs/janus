package oauth2

import (
	"context"
	"testing"
	"time"

	jwtbase "github.com/golang-jwt/jwt/v5"
	"github.com/hellofresh/janus/pkg/jwt"
	"github.com/hellofresh/janus/pkg/metrics"
	"github.com/hellofresh/stats-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJWTManagerValidKey(t *testing.T) {
	signingMethod := jwt.SigningMethod{Alg: "HS256", Key: "secret"}
	config := jwt.NewParserConfig(0, signingMethod)
	parser := jwt.NewParser(config)
	manager := NewJWTManager(parser)

	token, err := issueToken(signingMethod, 1*time.Hour)
	require.NoError(t, err)

	client, err := stats.NewClient("noop://")
	require.NoError(t, err)

	ctx := metrics.NewContext(context.Background(), client)
	assert.True(t, manager.IsKeyAuthorized(ctx, token))
}

func TestJWTManagerInvalidKey(t *testing.T) {
	signingMethod := jwt.SigningMethod{Alg: "HS256", Key: "secret"}
	config := jwt.NewParserConfig(0, signingMethod)
	parser := jwt.NewParser(config)
	manager := NewJWTManager(parser)

	client, err := stats.NewClient("noop://")
	require.NoError(t, err)

	ctx := metrics.NewContext(context.Background(), client)
	assert.False(t, manager.IsKeyAuthorized(ctx, "wrong"))
}

func TestJWTManagerNilContext(t *testing.T) {
	signingMethod := jwt.SigningMethod{Alg: "HS256", Key: "secret"}
	config := jwt.NewParserConfig(0, signingMethod)
	parser := jwt.NewParser(config)
	manager := NewJWTManager(parser)

	assert.False(t, manager.IsKeyAuthorized(nil, "wrong"))
}

func TestJWTManagerNilStast(t *testing.T) {
	signingMethod := jwt.SigningMethod{Alg: "HS256", Key: "secret"}
	config := jwt.NewParserConfig(0, signingMethod)
	parser := jwt.NewParser(config)
	manager := NewJWTManager(parser)

	assert.False(t, manager.IsKeyAuthorized(context.Background(), "wrong"))
}

func issueToken(signingMethod jwt.SigningMethod, expireIn time.Duration) (string, error) {
	token := jwtbase.New(jwtbase.GetSigningMethod(signingMethod.Alg))
	claims := token.Claims.(jwtbase.MapClaims)

	expire := time.Now().Add(expireIn)
	claims["exp"] = expire.Unix()
	claims["iat"] = time.Now().Unix()

	// currently only HSXXX algorithms are supported for issuing admin token, so we cast key to bytes array
	return token.SignedString([]byte(signingMethod.Key))
}
