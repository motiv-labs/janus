package middleware

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/hellofresh/janus/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// HostMatcher is a middleware that matches any host with the given list of hosts.
// It also supports regex host like *.example.com
type HostMatcher struct {
	plainHosts    map[string]bool
	wildcardHosts []*regexp.Regexp
}

// NewHostMatcher creates a new instance of HostMatcher
func NewHostMatcher(hosts []string) *HostMatcher {
	matcher := &HostMatcher{plainHosts: make(map[string]bool)}
	matcher.prepareIndexes(hosts)
	return matcher
}

// Handler is the middleware function
func (h *HostMatcher) Handler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithField("path", r.URL.Path).Debug("Starting host matcher middleware")
		host := r.Host

		if _, ok := h.plainHosts[host]; ok {
			log.WithField("host", host).Debug("Plain host matched")
			handler.ServeHTTP(w, r)
			return
		}

		for _, hostRegex := range h.wildcardHosts {
			if hostRegex.MatchString(host) {
				log.WithField("host", host).Debug("Wildcard host matched")
				handler.ServeHTTP(w, r)
				return
			}
		}

		err := errors.ErrRouteNotFound
		log.WithError(err).Error("The host didn't match any of the provided hosts")
		errors.Handler(w, r, err)
	})
}

func (h *HostMatcher) prepareIndexes(hosts []string) {
	if len(hosts) > 0 {
		for _, host := range hosts {
			if strings.Contains(host, "*") {
				regexStr := strings.Replace(host, ".", "\\.", -1)
				regexStr = strings.Replace(regexStr, "*", ".+", -1)
				h.wildcardHosts = append(h.wildcardHosts, regexp.MustCompile(fmt.Sprintf("^%s$", regexStr)))
			} else {
				h.plainHosts[host] = true
			}
		}
	}
}
