package jwt

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/hellofresh/janus/pkg/config"
	"github.com/hellofresh/janus/pkg/jwt/provider"
	"github.com/hellofresh/janus/pkg/render"
	"golang.org/x/oauth2"
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
func (j *Handler) Login(config config.Credentials) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// We're using OAuth, start checking for access keys
		authHeaderValue := r.Header.Get("Authorization")
		parts := strings.Split(authHeaderValue, " ")
		if len(parts) < 2 {
			render.JSON(w, http.StatusBadRequest, "attempted access with malformed header, no auth header found.")
			return
		}

		if strings.ToLower(parts[0]) != "bearer" {
			render.JSON(w, http.StatusBadRequest, "bearer token malformed")
			return
		}

		factory := provider.Factory{}
		p := factory.Build(config)

		ctx := context.Background()
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: parts[1]},
		)
		httpClient := oauth2.NewClient(ctx, ts)

		verified, err := p.Verify(httpClient)
		if err != nil {
			render.JSON(w, http.StatusInternalServerError, "failed to verify token")
			return
		}

		if !verified {
			render.JSON(w, http.StatusUnauthorized, "verification failed")
			return
		}

		if 0 == j.Guard.Timeout {
			j.Guard.Timeout = time.Hour
		}

		expire := time.Now().Add(j.Guard.Timeout)
		tokenString, err := IssueAdminToken(j.Guard.SigningMethod, map[string]interface{}{}, j.Guard.Timeout)
		if err != nil {
			render.JSON(w, http.StatusUnauthorized, "problem issuing JWT")
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
			render.JSON(w, http.StatusUnauthorized, "token is expired")
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
			render.JSON(w, http.StatusUnauthorized, "create JWT Token failed")
			return
		}

		render.JSON(w, http.StatusOK, render.M{
			"token":  tokenString,
			"expire": expire.Format(time.RFC3339),
		})
	}
}
