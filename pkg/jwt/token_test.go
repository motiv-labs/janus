package jwt

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIssueAdminToken(t *testing.T) {
	alg := "HS256"
	key := time.Now().Format(time.RFC3339Nano)
	claimsID := time.Now().Format(time.RFC3339Nano)

	accessToken, err := IssueAdminToken(SigningMethod{alg, key}, jwt.MapClaims{"id": claimsID}, time.Hour)
	require.NoError(t, err)

	config := NewParserConfig(0, SigningMethod{Alg: alg, Key: key})
	parser := NewParser(config)

	token, err := parser.Parse(accessToken.Token)
	require.NoError(t, err)

	claims, ok := parser.GetMapClaims(token)
	assert.True(t, ok)
	assert.Equal(t, claimsID, claims["id"])
}
