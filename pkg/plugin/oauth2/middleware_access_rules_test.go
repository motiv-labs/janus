package oauth2

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	basejwt "github.com/dgrijalva/jwt-go"
	"github.com/hellofresh/janus/pkg/jwt"
	"github.com/hellofresh/janus/pkg/oauth"
	"github.com/hellofresh/janus/pkg/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBlockJWTByCountry(t *testing.T) {
	secret := "secret"

	revokeRules := []*oauth.AccessRule{
		{Predicate: "country == 'de'", Action: "deny"},
	}

	config := jwt.NewConfig(secret)
	parser := jwt.NewParser(config)

	mw := NewRevokeRulesMiddleware(parser, revokeRules)
	token, err := generateToken(secret)
	require.NoError(t, err)

	w, err := test.Record(
		"GET",
		"/",
		map[string]string{
			"Content-Type":  "application/json",
			"Authorization": fmt.Sprintf("Bearer %s", token),
		},
		mw(http.HandlerFunc(test.Ping)),
	)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestBlockJWTByUsername(t *testing.T) {
	secret := "secret"

	revokeRules := []*oauth.AccessRule{
		{Predicate: "username == 'test@hellofresh.com'", Action: "deny"},
	}

	config := jwt.NewConfig(secret)
	parser := jwt.NewParser(config)

	mw := NewRevokeRulesMiddleware(parser, revokeRules)
	token, err := generateToken(secret)
	require.NoError(t, err)

	w, err := test.Record(
		"GET",
		"/",
		map[string]string{
			"Content-Type":  "application/json",
			"Authorization": fmt.Sprintf("Bearer %s", token),
		},
		mw(http.HandlerFunc(test.Ping)),
	)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestBlockJWTByIssueDate(t *testing.T) {
	secret := "secret"

	revokeRules := []*oauth.AccessRule{
		{Predicate: fmt.Sprintf("iat < %d", time.Now().Add(1*time.Hour).Unix()), Action: "deny"},
	}

	config := jwt.NewConfig(secret)
	parser := jwt.NewParser(config)

	mw := NewRevokeRulesMiddleware(parser, revokeRules)
	token, err := generateToken(secret)
	require.NoError(t, err)

	w, err := test.Record(
		"GET",
		"/",
		map[string]string{
			"Content-Type":  "application/json",
			"Authorization": fmt.Sprintf("Bearer %s", token),
		},
		mw(http.HandlerFunc(test.Ping)),
	)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestBlockJWTByCountryAndIssueDate(t *testing.T) {
	secret := "secret"

	revokeRules := []*oauth.AccessRule{
		{Predicate: fmt.Sprintf("country == 'de' && iat < %d", time.Now().Add(1*time.Hour).Unix()), Action: "deny"},
	}

	config := jwt.NewConfig(secret)
	parser := jwt.NewParser(config)

	mw := NewRevokeRulesMiddleware(parser, revokeRules)
	token, err := generateToken(secret)
	require.NoError(t, err)

	w, err := test.Record(
		"GET",
		"/",
		map[string]string{
			"Content-Type":  "application/json",
			"Authorization": fmt.Sprintf("Bearer %s", token),
		},
		mw(http.HandlerFunc(test.Ping)),
	)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func generateToken(secret string) (string, error) {
	token := basejwt.NewWithClaims(basejwt.SigningMethodHS256, basejwt.MapClaims{
		"country":  "de",
		"username": "test@hellofresh.com",
		"iat":      time.Now().Unix(),
	})

	return token.SignedString([]byte(secret))
}
