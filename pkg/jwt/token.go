package jwt

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

// IssueAdminToken issues admin JWT for API access
func IssueAdminToken(signingMethod SigningMethod, claims map[string]interface{}, expireIn time.Duration) (string, error) {
	token := jwt.New(jwt.GetSigningMethod(signingMethod.Alg))
	claims["exp"] = time.Now().Add(expireIn).Unix()
	claims["iat"] = time.Now().Unix()

	token.Claims = jwt.MapClaims(claims)

	// currently only HSXXX algorithms are supported for issuing admin token, so we cast key to bytes array
	return token.SignedString([]byte(signingMethod.Key))
}
