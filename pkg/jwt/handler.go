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
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

const (
	bearer = "bearer"
)

// Handler struct
type Handler struct {
	Guard Guard
}

// Login can be used by clients to get a jwt token.
// Payload needs to be json in the form of {"username": "<USERNAME>", "password": "<PASSWORD>"}.
// Reply will be of the form {"token": "<TOKEN>"}.
func (j *Handler) Login(config config.Credentials) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accessToken, err := extractAccessToken(r)
		if err != nil {
			log.WithError(err).Debug("failed to extract access token")
		}

		httpClient := getClient(accessToken)
		factory := provider.Factory{}
		p := factory.Build(r.URL.Query().Get("provider"), config)

		verified, err := p.Verify(r, httpClient)
		if err != nil || !verified {
			log.WithError(err).Debug(err.Error())
			render.JSON(w, http.StatusUnauthorized, err.Error())
			return
		}

		if 0 == j.Guard.Timeout {
			j.Guard.Timeout = time.Hour
		}

		claims, err := p.GetClaims(httpClient)
		if err != nil {
			render.JSON(w, http.StatusBadRequest, err.Error())
			return
		}

		token, err := IssueAdminToken(j.Guard.SigningMethod, claims, j.Guard.Timeout)
		if err != nil {
			render.JSON(w, http.StatusUnauthorized, "problem issuing JWT")
			return
		}

		render.JSON(w, http.StatusOK, token)
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
		newClaims["sub"] = claims["sub"]
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
			"type":   "Bearer",
			"expire": expire.Format(time.RFC3339),
		})
	}
}

func extractAccessToken(r *http.Request) (string, error) {
	// We're using OAuth, start checking for access keys
	authHeaderValue := r.Header.Get("Authorization")
	parts := strings.Split(authHeaderValue, " ")
	if len(parts) < 2 {
		return "", errors.New("attempted access with malformed header, no auth header found")
	}

	if strings.ToLower(parts[0]) != bearer {
		return "", errors.New("bearer token malformed")
	}

	return parts[1], nil
}

func getClient(token string) *http.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	return oauth2.NewClient(ctx, ts)
}
