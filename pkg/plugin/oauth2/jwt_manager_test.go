package oauth2

import (
	"testing"
	"time"

	jwtbase "github.com/dgrijalva/jwt-go"
	"github.com/hellofresh/janus/pkg/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJWTManagerValidKey(t *testing.T) {
	signingMethod := jwt.SigningMethod{Alg: "HS256", Key: "secret"}
	config := jwt.NewParserConfig(signingMethod)
	parser := jwt.NewParser(config)
	manager := NewJWTManager(parser)

	token, err := issueToken(signingMethod, 1*time.Hour)
	require.NoError(t, err)

	assert.True(t, manager.IsKeyAuthorized(token))
}

func TestJWTManagerInvalidKey(t *testing.T) {
	signingMethod := jwt.SigningMethod{Alg: "HS256", Key: "secret"}
	config := jwt.NewParserConfig(signingMethod)
	parser := jwt.NewParser(config)
	manager := NewJWTManager(parser)

	assert.False(t, manager.IsKeyAuthorized("wrong"))
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
