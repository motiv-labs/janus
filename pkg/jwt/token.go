package jwt

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// IssueAdminToken issues admin JWT for API access
func IssueAdminToken(signingAlgorithm, claimsID string, secret []byte, expireIn time.Duration) (string, error) {
	token := jwt.New(jwt.GetSigningMethod(signingAlgorithm))
	claims := token.Claims.(jwt.MapClaims)

	expire := time.Now().Add(expireIn)
	claims["id"] = claimsID
	claims["exp"] = expire.Unix()
	claims["iat"] = time.Now().Unix()

	return token.SignedString(secret)
}
