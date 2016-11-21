package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/appleboy/gin-jwt.v2"
)

type Credentials struct {
	Secret, Username, Password string
}

func NewJwt(cred *Credentials) *jwt.GinJWTMiddleware {
	return &jwt.GinJWTMiddleware{
		Realm:      "test zone",
		Key:        []byte(cred.Secret),
		Timeout:    time.Hour,
		MaxRefresh: time.Hour * 24,
		Authenticator: func(userId string, password string, c *gin.Context) (string, bool) {
			if userId == cred.Username && password == cred.Password {
				return userId, true
			}

			return userId, false
		},
		Authorizator: func(userId string, c *gin.Context) bool {
			if userId == cred.Username {
				return true
			}

			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
	}
}
