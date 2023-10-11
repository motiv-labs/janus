package authorization

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	jwtUser "github.com/golang-jwt/jwt/v4"
)

type UserClaims struct {
	UserID       int
	FullUserName string
	Roles        []string
	jwtUser.StandardClaims
}

func ExtractClaims(jwtToken string) (*UserClaims, error) {
	parts := strings.Split(jwtToken, ".")
	if len(parts) < 2 {
		return nil, fmt.Errorf("JWT token invalid")
	}

	decoded, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}

	claims := &UserClaims{}
	err = json.Unmarshal([]byte(decoded), claims)
	if err != nil {
		return nil, err
	}
	return claims, nil
}
