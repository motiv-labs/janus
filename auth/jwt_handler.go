package auth

import (
	"errors"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/hellofresh/janus/request"
	"github.com/hellofresh/janus/response"
)

// JWTHandler struct
type JWTHandler struct {
	Config JWTConfig
}

// Login form structure.
type Login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// Login can be used by clients to get a jwt token.
// Payload needs to be json in the form of {"username": "USERNAME", "password": "PASSWORD"}.
// Reply will be of the form {"token": "TOKEN"}.
func (j *JWTHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var loginVals Login

		if request.BindJSON(r, &loginVals) != nil {
			j.Config.Unauthorized(w, r, errors.New("missing username or password"))
			return
		}

		userID, ok := j.Config.Authenticator(loginVals.Username, loginVals.Password)

		if !ok {
			j.Config.Unauthorized(w, r, errors.New("invalid username or password"))
			return
		}

		// Create the token
		token := jwt.New(jwt.GetSigningMethod(j.Config.SigningAlgorithm))
		claims := token.Claims.(jwt.MapClaims)

		if userID == "" {
			userID = loginVals.Username
		}

		if 0 == j.Config.Timeout {
			j.Config.Timeout = time.Hour
		}

		expire := time.Now().Add(j.Config.Timeout)
		claims["id"] = userID
		claims["exp"] = expire.Unix()
		claims["iat"] = time.Now().Unix()

		tokenString, err := token.SignedString(j.Config.Secret)

		if err != nil {
			j.Config.Unauthorized(w, r, errors.New("problem signing JWT"))
			return
		}

		response.JSON(w, http.StatusOK, response.H{
			"token":  tokenString,
			"expire": expire.Format(time.RFC3339),
		})
	}
}

func (j *JWTHandler) Refresh() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		parser := JWTParser{j.Config}
		token, _ := parser.Parse(r)
		claims := token.Claims.(jwt.MapClaims)

		origIat := int64(claims["iat"].(float64))

		if origIat < time.Now().Add(-j.Config.MaxRefresh).Unix() {
			j.Config.Unauthorized(w, r, errors.New("token is expired"))
			return
		}

		// Create the token
		newToken := jwt.New(jwt.GetSigningMethod(j.Config.SigningAlgorithm))
		newClaims := newToken.Claims.(jwt.MapClaims)

		for key := range claims {
			newClaims[key] = claims[key]
		}

		expire := time.Now().Add(j.Config.Timeout)
		newClaims["id"] = claims["id"]
		newClaims["exp"] = expire.Unix()
		newClaims["orig_iat"] = origIat

		tokenString, err := newToken.SignedString(j.Config.Secret)

		if err != nil {
			j.Config.Unauthorized(w, r, errors.New("create JWT Token faild"))
			return
		}

		response.JSON(w, http.StatusOK, response.H{
			"token":  tokenString,
			"expire": expire.Format(time.RFC3339),
		})
	}
}
