package jwt

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

// IssueAdminToken issues admin JWT for API access
func IssueAdminToken(signingMethod SigningMethod, claimsID string, expireIn time.Duration) (string, error) {
	token := jwt.New(jwt.GetSigningMethod(signingMethod.Alg))
	claims := token.Claims.(jwt.MapClaims)

	expire := time.Now().Add(expireIn)
	claims["id"] = claimsID
	claims["exp"] = expire.Unix()
	claims["iat"] = time.Now().Unix()

	// currently only HSXXX algorithms are supported for issuing admin token, so we cast key to bytes array
	return token.SignedString([]byte(signingMethod.Key))
}
