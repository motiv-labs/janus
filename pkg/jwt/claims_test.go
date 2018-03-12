package jwt

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJanusClaims_UnmarshalJSON(t *testing.T) {
	claims := NewJanusClaims(0)

	claimBytes := []byte(`{"exp":1,"iat":2,"iss":"test","username":"janus"}`)
	dec := json.NewDecoder(bytes.NewBuffer(claimBytes))
	err := dec.Decode(&claims)
	require.NoError(t, err)

	assert.Equal(t, float64(1), claims.MapClaims["exp"])
	assert.Equal(t, float64(2), claims.MapClaims["iat"])
	assert.Equal(t, "test", claims.MapClaims["iss"])
	assert.Equal(t, "janus", claims.MapClaims["username"])
}

func TestJanusClaims_VerifyExpiresAt(t *testing.T) {
	leeway := 1 + rand.Int63n(120)
	claims := NewJanusClaims(leeway)
	now := time.Now().Unix()

	claims.MapClaims["exp"] = float64(now - 1)
	assert.True(t, claims.VerifyExpiresAt(now, true))

	claims.MapClaims["exp"] = float64(now - leeway)
	assert.True(t, claims.VerifyExpiresAt(now, true))

	claims.MapClaims["exp"] = float64(now - leeway - 1)
	assert.False(t, claims.VerifyExpiresAt(now, true))
}

func TestJanusClaims_VerifyIssuedAt(t *testing.T) {
	leeway := 1 + rand.Int63n(120)
	claims := NewJanusClaims(leeway)
	now := time.Now().Unix()

	claims.MapClaims["iat"] = float64(now + 1)
	assert.True(t, claims.VerifyIssuedAt(now, true))

	claims.MapClaims["iat"] = float64(now + leeway)
	assert.True(t, claims.VerifyIssuedAt(now, true))

	claims.MapClaims["iat"] = float64(now + leeway + 1)
	assert.False(t, claims.VerifyIssuedAt(now, true))
}

func TestJanusClaims_VerifyNotBefore(t *testing.T) {
	leeway := 1 + rand.Int63n(120)
	claims := NewJanusClaims(leeway)
	now := time.Now().Unix()

	claims.MapClaims["nbf"] = float64(now + 1)
	assert.True(t, claims.VerifyNotBefore(now, true))

	claims.MapClaims["nbf"] = float64(now + leeway)
	assert.True(t, claims.VerifyNotBefore(now, true))

	claims.MapClaims["nbf"] = float64(now + leeway + 1)
	assert.False(t, claims.VerifyNotBefore(now, true))
}
