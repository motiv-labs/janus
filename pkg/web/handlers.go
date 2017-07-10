package web

import (
	"fmt"
	"net"
	"net/http"
	"net/url"

	"github.com/hellofresh/janus/pkg/errors"
	"github.com/hellofresh/janus/pkg/response"
	log "github.com/sirupsen/logrus"
)

// Home handler is just a nice home page message
func Home(version string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response.JSON(w, http.StatusOK, "Welcome to Janus v"+version)
	}
}

// NotFound handler is called when no route is matched
func NotFound(w http.ResponseWriter, r *http.Request) {
	notFoundError := errors.ErrRouteNotFound
	response.JSON(w, notFoundError.Code, notFoundError)
}

// RecoveryHandler handler is used when a panic happens
func RecoveryHandler(w http.ResponseWriter, r *http.Request, err interface{}) {
	switch internalErr := err.(type) {
	case *errors.Error:
		log.WithFields(log.Fields{"code": internalErr.Code, "error": internalErr.Error()}).
			Info("Internal error handled")
		response.JSON(w, internalErr.Code, internalErr.Error())
	default:
		log.WithField("error", err).Error("Internal server error handled")
		response.JSON(w, http.StatusInternalServerError, err)
	}
}

// RedirectHTTPS redirects an http request to https
func RedirectHTTPS(port int) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		host, _, _ := net.SplitHostPort(req.Host)

		target := url.URL{
			Scheme: "https",
			Host:   fmt.Sprintf("%s:%v", host, port),
			Path:   req.URL.Path,
		}
		if len(req.URL.RawQuery) > 0 {
			target.RawQuery += "?" + req.URL.RawQuery
		}
		log.Printf("redirect to: %s", target.String())
		http.Redirect(w, req, target.String(), http.StatusTemporaryRedirect)
	})
}
