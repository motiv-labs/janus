package jwt

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIssueAdminToken(t *testing.T) {
	alg := "HS256"
	key := time.Now().Format(time.RFC3339Nano)
	claimsID := time.Now().Format(time.RFC3339Nano)

	tokenString, err := IssueAdminToken(SigningMethod{alg, key}, claimsID, time.Hour)
	require.NoError(t, err)

	config := NewParserConfig(SigningMethod{Alg: alg, Key: key})
	parser := NewParser(config)

	token, err := parser.Parse(tokenString)
	require.NoError(t, err)

	claims, ok := parser.GetMapClaims(token)
	assert.True(t, ok)
	assert.Equal(t, claimsID, claims["id"])
}
