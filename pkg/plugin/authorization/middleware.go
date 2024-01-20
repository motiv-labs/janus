package authorization

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strings"

	"github.com/hellofresh/janus/pkg/errors"
	"github.com/hellofresh/janus/pkg/render"
)

var (
	AuthHeaderValue = ContextKey("auth_header")

	ErrAuthorizationFieldNotFound = errors.New(http.StatusBadRequest, "authorization field missing")
	ErrBearerMalformed            = errors.New(http.StatusBadRequest, "bearer token malformed")
	ErrAccessTokenNotAuthorized   = errors.New(http.StatusUnauthorized, "access token not authorized")
	ErrNoRolesSet                 = errors.New(http.StatusUnauthorized, "no roles in access token")
	ErrAccessIsDenied             = errors.New(http.StatusUnauthorized, "access is denied")
	ErrBodyReading                = errors.New(http.StatusInternalServerError, "body reading error")
	ErrUnmarshal                  = errors.New(http.StatusInternalServerError, "cannot unmarshal")
)

// ContextKey is used to create context keys that are concurrent safe
type ContextKey string

func (c ContextKey) String() string {
	return "janus." + string(c)
}

func NewTokenCheckerMiddleware(manager *TokenManager) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeaderValue := r.Header.Get("Authorization")
			parts := strings.Split(authHeaderValue, " ")

			if len(parts) == 0 {
				errors.Handler(w, r, ErrAuthorizationFieldNotFound)
				return
			} else if len(parts) < 2 {
				logrus.Errorf("bearer token malformed, token is: %q", authHeaderValue)
				errors.Handler(w, r, ErrBearerMalformed)
				return
			}

			accessToken := parts[1]

			err := manager.FetchTokens()
			if err != nil {
				errors.Handler(w, r, errors.New(http.StatusInternalServerError, err.Error()))
				return
			}

			if !isTokenAuthorized(manager.Tokens, accessToken) {
				errors.Handler(w, r, ErrAccessTokenNotAuthorized)
				return
			}

			ctx := context.WithValue(r.Context(), AuthHeaderValue, accessToken)
			handler.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func isTokenAuthorized(tokens map[string]*JWTToken, userToken string) bool {
	if _, exists := tokens[userToken]; exists {
		return true
	}

	return false
}

func NewRoleCheckerMiddleware(manager *RoleManager) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeaderValue := r.Header.Get("Authorization")
			parts := strings.Split(authHeaderValue, " ")

			if len(parts) == 0 {
				errors.Handler(w, r, ErrAuthorizationFieldNotFound)
				return
			} else if len(parts) < 2 {
				logrus.Errorf("bearer token malformed, token is: %q", authHeaderValue)
				errors.Handler(w, r, ErrBearerMalformed)
				return
			}

			accessToken := parts[1]

			claims, err := ExtractClaims(accessToken)
			if err != nil {
				errors.Handler(w, r, err)
				return
			}

			if len(claims.Roles) == 0 {
				errors.Handler(w, r, ErrNoRolesSet)
				return
			}

			err = manager.FetchRoles()
			if err != nil {
				errors.Handler(w, r, errors.New(http.StatusInternalServerError, err.Error()))
				return
			}

			if !isHaveAccess(manager.Roles, claims.Roles, r.URL.Path, r.Method) {
				errors.Handler(w, r, ErrAccessIsDenied)
				return
			}

			ctx := context.WithValue(r.Context(), AuthHeaderValue, accessToken)
			handler.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func isHaveAccess(roles map[string]*Role, userRoles []string, path, method string) bool {
	for _, userRole := range userRoles {
		if role, exists := roles[userRole]; exists {
			for _, feature := range role.Features {
				if feature.Method == method && isEndpointPathsEqual(path, feature.Path) {
					return true
				}
			}
		}
	}

	return false
}

func isEndpointPathsEqual(reqPath, confPath string) bool {
	reqPathArr := strings.Split(confPath, "/")
	confPathArr := strings.Split(reqPath, "/")
	if len(reqPathArr) != len(confPathArr) {
		return false
	}

	for i := range confPathArr {
		if reqPathArr[i] == "" || string(reqPathArr[i][0]) == "{" {
			continue
		}

		if reqPathArr[i] != confPathArr[i] {
			return false
		}
	}

	return true
}

func NewTokenCatcherMiddleware(manager *TokenManager) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case fmt.Sprintf("/%s/gatewayTokens", manager.Conf.ApiVersion):
				body, err := io.ReadAll(r.Body)
				if err != nil {
					errors.Handler(w, r, ErrBodyReading)
				}

				tokens := []*JWTToken{}
				err = json.Unmarshal(body, &tokens)
				if err != nil {
					errors.Handler(w, r, ErrUnmarshal)
				}

				manager.UpsertTokens(tokens)

				render.JSON(w, http.StatusOK, http.NoBody)
				return

			case fmt.Sprintf("/%s/gatewayTokens/delete", manager.Conf.ApiVersion):
				body, err := io.ReadAll(r.Body)
				if err != nil {
					errors.Handler(w, r, ErrBodyReading)
				}

				ids := []uint64{}
				err = json.Unmarshal(body, &ids)
				if err != nil {
					errors.Handler(w, r, ErrUnmarshal)
				}

				manager.DeleteTokensByIDs(ids)

				render.JSON(w, http.StatusOK, http.NoBody)
				return

			default:
				handler.ServeHTTP(w, r)
				return
			}
		})
	}
}

func NewRoleCatcherMiddleware(manager *RoleManager) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case fmt.Sprintf("/%s/gatewayRoles", manager.Conf.ApiVersion):
				body, err := io.ReadAll(r.Body)
				if err != nil {
					errors.Handler(w, r, ErrBodyReading)
				}

				roles := []*Role{}
				err = json.Unmarshal(body, &roles)
				if err != nil {
					errors.Handler(w, r, ErrUnmarshal)
				}

				manager.UpsertRoles(roles)

				render.JSON(w, http.StatusOK, http.NoBody)
				return

			case fmt.Sprintf("/%s/gatewayRoles/delete", manager.Conf.ApiVersion):
				body, err := io.ReadAll(r.Body)
				if err != nil {
					errors.Handler(w, r, ErrBodyReading)
				}

				ids := []uint64{}
				err = json.Unmarshal(body, &ids)
				if err != nil {
					errors.Handler(w, r, ErrUnmarshal)
				}

				manager.DeleteRolesByIDs(ids)

				render.JSON(w, http.StatusOK, http.NoBody)
				return

			default:
				handler.ServeHTTP(w, r)
				return
			}

		})
	}
}
