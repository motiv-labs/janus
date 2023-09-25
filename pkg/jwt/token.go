package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// AccessToken represents a token
type AccessToken struct {
	Type    string `json:"token_type"`
	Token   string `json:"access_token"`
	Expires int64  `json:"expires_in"`
}

// IssueAdminToken issues admin JWT for API access
func IssueAdminToken(signingMethod SigningMethod, claims jwt.MapClaims, expireIn time.Duration) (*AccessToken, error) {
	token := jwt.New(jwt.GetSigningMethod(signingMethod.Alg))
	exp := time.Now().Add(expireIn).Unix()

	token.Claims = claims
	claims["exp"] = exp
	claims["iat"] = time.Now().Unix()

	accessToken, err := token.SignedString([]byte(signingMethod.Key))
	if err != nil {
		return nil, err
	}

	// currently only HSXXX algorithms are supported for issuing admin token, so we cast key to bytes array
	return &AccessToken{
		Type:    "Bearer",
		Token:   accessToken,
		Expires: exp,
	}, nil
}
