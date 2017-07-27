package web

import (
	"fmt"
	"net"
	"net/http"
	"net/url"

	"github.com/hellofresh/janus/pkg/render"
	log "github.com/sirupsen/logrus"
)

// Home handler is just a nice home page message
func Home(version string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, http.StatusOK, "Welcome to Janus v"+version)
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
