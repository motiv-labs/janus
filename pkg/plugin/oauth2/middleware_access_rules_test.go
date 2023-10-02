package oauth2

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	basejwt "github.com/dgrijalva/jwt-go"
	"github.com/hellofresh/janus/pkg/jwt"
	"github.com/hellofresh/janus/pkg/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const signingAlg = "HS256"

func generateToken(alg, key string) (string, error) {
	token := basejwt.NewWithClaims(basejwt.GetSigningMethod(alg), basejwt.MapClaims{
		"country":  "de",
		"username": "test@hellofresh.com",
		"iat":      time.Now().Unix(),
	})

	return token.SignedString([]byte(key))
}

func expectRulesToProduceStatus(t *testing.T, statusCode int, rules []*AccessRule) {
	secret := "secret"

	parser := jwt.NewParser(jwt.NewParserConfig(0, jwt.SigningMethod{Alg: signingAlg, Key: secret}))

	mw := NewRevokeRulesMiddleware(parser, rules)
	token, err := generateToken(signingAlg, secret)
	require.NoError(t, err)

	for i := 1; i <= 3; i++ { // middleware caches predicate and should return the same response every time
		hits := 0
		w, err := test.Record(
			"GET",
			"/",
			map[string]string{
				"Content-Type":  "application/json",
				"Authorization": fmt.Sprintf("Bearer %s", token),
			},
			mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				hits++
				test.Ping(w, r)
			})),
		)

		assert.NoError(t, err, "%d. pass", i)
		assert.Equal(t, statusCode, w.Code, "%d. pass", i)
		if statusCode == http.StatusOK {
			assert.Equal(t, 1, hits, "%d. pass", i)
		} else {
			assert.Equal(t, 0, hits, "%d. pass", i)
		}
	}
}

func TestBlockJWTByCountry(t *testing.T) {
	expectRulesToProduceStatus(t, http.StatusUnauthorized, []*AccessRule{
		{Predicate: "country == 'de'", Action: "deny"},
	})
}

func TestBlockJWTByUsername(t *testing.T) {
	expectRulesToProduceStatus(t, http.StatusUnauthorized, []*AccessRule{
		{Predicate: "username == 'test@hellofresh.com'", Action: "deny"},
	})
}

func TestBlockJWTByIssueDate(t *testing.T) {
	expectRulesToProduceStatus(t, http.StatusUnauthorized, []*AccessRule{
		{Predicate: fmt.Sprintf("iat < %d", time.Now().Add(1*time.Hour).Unix()), Action: "deny"},
	})
}

func TestBlockJWTByCountryAndIssueDate(t *testing.T) {
	expectRulesToProduceStatus(t, http.StatusUnauthorized, []*AccessRule{
		{Predicate: fmt.Sprintf("country == 'de' && iat < %d", time.Now().Add(1*time.Hour).Unix()), Action: "deny"},
	})
}

func TestEmptyAccessRules(t *testing.T) {
	expectRulesToProduceStatus(t, http.StatusOK, []*AccessRule{})
}

func TestWrongRule(t *testing.T) {
	expectRulesToProduceStatus(t, http.StatusOK, []*AccessRule{
		{Predicate: "country == 'wrong'", Action: "deny"},
	})
}

func TestMultipleRulesSecondMatchesAndDenies(t *testing.T) {
	expectRulesToProduceStatus(t, http.StatusUnauthorized, []*AccessRule{
		{Predicate: "country == 'us'", Action: "deny"},
		{Predicate: "country == 'de'", Action: "deny"},
	})
}

func TestMultipleRulesSecondMatchesAndAllows(t *testing.T) {
	expectRulesToProduceStatus(t, http.StatusOK, []*AccessRule{
		{Predicate: "country == 'us'", Action: "allow"},
		{Predicate: "country == 'de'", Action: "allow"},
		{Predicate: "true", Action: "deny"},
	})
}

func TestMultipleRulesLastMatchesAndDenies(t *testing.T) {
	expectRulesToProduceStatus(t, http.StatusUnauthorized, []*AccessRule{
		{Predicate: "country == 'us'", Action: "allow"},
		{Predicate: "country == 'gb'", Action: "allow"},
		{Predicate: "true", Action: "deny"},
	})
}

func TestMultipleRulesNoneMatch(t *testing.T) {
	expectRulesToProduceStatus(t, http.StatusOK, []*AccessRule{
		{Predicate: "country == 'us'", Action: "deny"},
		{Predicate: "country == 'gb'", Action: "deny"},
	})
}
func TestMultipleRulesMatchAndAllow(t *testing.T) {
	expectRulesToProduceStatus(t, http.StatusOK, []*AccessRule{
		{Predicate: "country == 'de'", Action: "allow"},
		{Predicate: "true", Action: "allow"},
	})
}

func TestWrongJWT(t *testing.T) {
	revokeRules := []*AccessRule{
		{Predicate: fmt.Sprintf("country == 'de' && iat < %d", time.Now().Add(1*time.Hour).Unix()), Action: "deny"},
	}

	parser := jwt.NewParser(jwt.NewParserConfig(0, jwt.SigningMethod{Alg: signingAlg, Key: "secret"}))

	mw := NewRevokeRulesMiddleware(parser, revokeRules)
	token, err := generateToken(signingAlg, "wrong secret")
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
