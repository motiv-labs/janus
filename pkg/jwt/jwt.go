package jwt

import "github.com/golang-jwt/jwt/v4"

type Claims struct {
	UserID       int
	FullUserName string
	jwt.StandardClaims
}
