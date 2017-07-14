package jwt

import (
	"errors"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/hellofresh/janus/pkg/request"
	"github.com/hellofresh/janus/pkg/response"
)

// Handler struct
type Handler struct {
	Config Config
}

// Login form structure.
type Login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// Login can be used by clients to get a jwt token.
// Payload needs to be json in the form of {"username": "<USERNAME>", "password": "<PASSWORD>"}.
// Reply will be of the form {"token": "<TOKEN>"}.
func (j *Handler) Login() http.HandlerFunc {
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

		if userID == "" {
			userID = loginVals.Username
		}

		if 0 == j.Config.Timeout {
			j.Config.Timeout = time.Hour
		}

		expire := time.Now().Add(j.Config.Timeout)

		tokenString, err := IssueAdminToken(j.Config.SigningAlgorithm, userID, j.Config.Secret, j.Config.Timeout)

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

// Refresh can be used to refresh existing and valid jwt token.
// Reply will be of the form {"token": "<TOKEN>", "expire": "<DateTime in RFC-3339 format>"}.
func (j *Handler) Refresh() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		parser := Parser{j.Config}
		token, _ := parser.ParseFromRequest(r)
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
		newClaims["iat"] = origIat

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
