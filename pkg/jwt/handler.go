package jwt

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/hellofresh/janus/pkg/render"
)

// Handler struct
type Handler struct {
	Guard Guard
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
		var loginValues Login

		if json.NewDecoder(r.Body).Decode(&loginValues) != nil {
			j.Guard.Unauthorized(w, r, errors.New("missing username or password"))
			return
		}

		userID, ok := j.Guard.Authenticator(loginValues.Username, loginValues.Password)

		if !ok {
			j.Guard.Unauthorized(w, r, errors.New("invalid username or password"))
			return
		}

		if userID == "" {
			userID = loginValues.Username
		}

		if 0 == j.Guard.Timeout {
			j.Guard.Timeout = time.Hour
		}

		expire := time.Now().Add(j.Guard.Timeout)

		tokenString, err := IssueAdminToken(j.Guard.SigningMethod, userID, j.Guard.Timeout)

		if err != nil {
			j.Guard.Unauthorized(w, r, errors.New("problem issuing JWT"))
			return
		}

		render.JSON(w, http.StatusOK, render.M{
			"token":  tokenString,
			"expire": expire.Format(time.RFC3339),
		})
	}
}

// Refresh can be used to refresh existing and valid jwt token.
// Reply will be of the form {"token": "<TOKEN>", "expire": "<DateTime in RFC-3339 format>"}.
func (j *Handler) Refresh() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		parser := Parser{j.Guard.ParserConfig}
		token, _ := parser.ParseFromRequest(r)
		claims := token.Claims.(jwt.MapClaims)

		origIat := int64(claims["iat"].(float64))

		if origIat < time.Now().Add(-j.Guard.MaxRefresh).Unix() {
			j.Guard.Unauthorized(w, r, errors.New("token is expired"))
			return
		}

		// Create the token
		newToken := jwt.New(jwt.GetSigningMethod(j.Guard.SigningMethod.Alg))
		newClaims := newToken.Claims.(jwt.MapClaims)

		for key := range claims {
			newClaims[key] = claims[key]
		}

		expire := time.Now().Add(j.Guard.Timeout)
		newClaims["id"] = claims["id"]
		newClaims["exp"] = expire.Unix()
		newClaims["iat"] = origIat

		// currently only HSXXX algorithms are supported for issuing admin token, so we cast key to bytes array
		tokenString, err := newToken.SignedString([]byte(j.Guard.SigningMethod.Key))
		if err != nil {
			j.Guard.Unauthorized(w, r, errors.New("create JWT Token failed"))
			return
		}

		render.JSON(w, http.StatusOK, render.M{
			"token":  tokenString,
			"expire": expire.Format(time.RFC3339),
		})
	}
}
