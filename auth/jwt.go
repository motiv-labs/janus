package auth

import (
	"net/http"
	"time"

	"github.com/hellofresh/janus/config"
	"github.com/hellofresh/janus/response"
)

func NewJWTConfig(cred config.Credentials) JWTConfig {
	return JWTConfig{
		SigningAlgorithm: "HS256",
		Secret:           []byte(cred.Secret),
		Timeout:          time.Hour,
		MaxRefresh:       time.Hour * 24,
		Authenticator: func(userID string, password string) (string, bool) {
			if userID == cred.Username && password == cred.Password {
				return userID, true
			}

			return userID, false
		},
		Authorizator: func(userID string, w http.ResponseWriter, r *http.Request) bool {
			if userID == cred.Username {
				return true
			}

			return false
		},
		Unauthorized: func(w http.ResponseWriter, r *http.Request, err error) {
			response.JSON(w, http.StatusUnauthorized, response.H{
				"message": err.Error(),
			})
		},
		TokenLookup: "header:Authorization",
	}
}
