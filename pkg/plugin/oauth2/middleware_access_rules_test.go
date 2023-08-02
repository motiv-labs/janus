package oauth2

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	basejwt "github.com/golang-jwt/jwt/v5"
	"github.com/hellofresh/janus/pkg/jwt"
	"github.com/hellofresh/janus/pkg/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const signingAlg = "HS256"

func TestBlockJWTByCountry(t *testing.T) {
	secret := "secret"

	revokeRules := []*AccessRule{
		{Predicate: "country == 'de'", Action: "deny"},
	}

	parser := jwt.NewParser(jwt.NewParserConfig(0, jwt.SigningMethod{Alg: signingAlg, Key: secret}))

	mw := NewRevokeRulesMiddleware(parser, revokeRules)
	token, err := generateToken(signingAlg, secret)
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

	revokeRules := []*AccessRule{
		{Predicate: "username == 'test@hellofresh.com'", Action: "deny"},
	}

	parser := jwt.NewParser(jwt.NewParserConfig(0, jwt.SigningMethod{Alg: signingAlg, Key: secret}))

	mw := NewRevokeRulesMiddleware(parser, revokeRules)
	token, err := generateToken(signingAlg, secret)
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

	revokeRules := []*AccessRule{
		{Predicate: fmt.Sprintf("iat < %d", time.Now().Add(1*time.Hour).Unix()), Action: "deny"},
	}

	parser := jwt.NewParser(jwt.NewParserConfig(0, jwt.SigningMethod{Alg: signingAlg, Key: secret}))

	mw := NewRevokeRulesMiddleware(parser, revokeRules)
	token, err := generateToken(signingAlg, secret)
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

	revokeRules := []*AccessRule{
		{Predicate: fmt.Sprintf("country == 'de' && iat < %d", time.Now().Add(1*time.Hour).Unix()), Action: "deny"},
	}

	parser := jwt.NewParser(jwt.NewParserConfig(0, jwt.SigningMethod{Alg: signingAlg, Key: secret}))

	mw := NewRevokeRulesMiddleware(parser, revokeRules)
	token, err := generateToken(signingAlg, secret)
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

func generateToken(alg, key string) (string, error) {
	token := basejwt.NewWithClaims(basejwt.GetSigningMethod(alg), basejwt.MapClaims{
		"country":  "de",
		"username": "test@hellofresh.com",
		"iat":      time.Now().Unix(),
	})

	return token.SignedString([]byte(key))
}

func TestEmptyAccessRules(t *testing.T) {
	secret := "secret"

	revokeRules := []*AccessRule{}

	parser := jwt.NewParser(jwt.NewParserConfig(0, jwt.SigningMethod{Alg: signingAlg, Key: secret}))

	mw := NewRevokeRulesMiddleware(parser, revokeRules)

	w, err := test.Record(
		"GET",
		"/",
		nil,
		mw(http.HandlerFunc(test.Ping)),
	)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestWrongJWT(t *testing.T) {
	revokeRules := []*AccessRule{
		{Predicate: fmt.Sprintf("country == 'de' && iat < %d", time.Now().Add(1*time.Hour).Unix()), Action: "deny"},
	}

	parser := jwt.NewParser(jwt.NewParserConfig(0, jwt.SigningMethod{Alg: signingAlg, Key: "wrong secret"}))

	mw := NewRevokeRulesMiddleware(parser, revokeRules)
	token, err := generateToken(signingAlg, "secret")
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
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestWrongRule(t *testing.T) {
	secret := "secret"

	revokeRules := []*AccessRule{
		{Predicate: "country == 'wrong'", Action: "deny"},
	}

	parser := jwt.NewParser(jwt.NewParserConfig(0, jwt.SigningMethod{Alg: signingAlg, Key: secret}))

	mw := NewRevokeRulesMiddleware(parser, revokeRules)
	token, err := generateToken(signingAlg, secret)
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
	assert.Equal(t, http.StatusOK, w.Code)
}
