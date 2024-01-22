package authorization

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	UserID       uint64
	FullUserName string
	jwt.StandardClaims
	Roles []string
}

func ExtractClaims(jwtToken string) (*Claims, error) {
	parts := strings.Split(jwtToken, ".")
	if len(parts) < 2 {
		return nil, fmt.Errorf("JWT token invalid")
	}

	decoded, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}

	claims := &Claims{}
	err = json.Unmarshal(decoded, claims)
	if err != nil {
		return nil, err
	}
	return claims, nil
}
