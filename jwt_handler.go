package janus

import (
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/hellofresh/janus/request"
	"github.com/hellofresh/janus/response"
)

// JWTConfig struct
type JWTConfig struct {
	Authenticator    func(userID string, password string) (string, bool)
	Timeout          time.Duration
	Secret           []byte
	SigningAlgorithm string
}

// Login form structure.
type Login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// Login can be used by clients to get a jwt token.
// Payload needs to be json in the form of {"username": "USERNAME", "password": "PASSWORD"}.
// Reply will be of the form {"token": "TOKEN"}.
func (j *JWTConfig) Login() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var loginVals Login

		if request.BindJSON(r, &loginVals) != nil {
			response.JSON(rw, http.StatusBadRequest, response.H{
				"error": "Missing username or password",
			})
			return
		}

		userID, ok := j.Authenticator(loginVals.Username, loginVals.Password)

		if !ok {
			response.JSON(rw, http.StatusUnauthorized, response.H{
				"error": "Invalid username or password",
			})
			return
		}

		// Create the token
		token := jwt.New(jwt.GetSigningMethod(j.SigningAlgorithm))
		claims := token.Claims.(jwt.MapClaims)

		if userID == "" {
			userID = loginVals.Username
		}

		if 0 == j.Timeout {
			j.Timeout = time.Hour
		}

		expire := time.Now().Add(j.Timeout)
		claims["id"] = userID
		claims["exp"] = expire.Unix()
		claims["iat"] = time.Now().Unix()

		tokenString, err := token.SignedString(j.Secret)

		if err != nil {
			response.JSON(rw, http.StatusUnauthorized, response.H{
				"error": "Problem signing JWT",
			})
			return
		}

		response.JSON(rw, http.StatusOK, response.H{
			"token":  tokenString,
			"expire": expire.Format(time.RFC3339),
		})
	}
}
