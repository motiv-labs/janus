package authorization

import (
	"errors"
	"net/http"

	errorsJanus "github.com/hellofresh/janus/pkg/errors"
)

var (
	ErrAuthorizationFieldNotFound = errorsJanus.New(http.StatusBadRequest, "authorization field missing")
	ErrBearerMalformed            = errorsJanus.New(http.StatusBadRequest, "bearer token malformed")
	ErrAccessTokenNotAuthorized   = errorsJanus.New(http.StatusUnauthorized, "access token not authorized")
	ErrNoRolesSet                 = errorsJanus.New(http.StatusUnauthorized, "no roles in access token")
	ErrAccessIsDenied             = errorsJanus.New(http.StatusUnauthorized, "access is denied")
	ErrBodyReading                = errorsJanus.New(http.StatusInternalServerError, "body reading error")
	ErrUnmarshal                  = errorsJanus.New(http.StatusInternalServerError, "cannot unmarshal")
)

var (
	ErrEventTypeConvert = errors.New("cannot convert event")
)
