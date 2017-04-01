package plugin

import (
	"regexp"

	"github.com/hellofresh/janus/pkg/api"
	"github.com/hellofresh/janus/pkg/middleware"
	"github.com/hellofresh/janus/pkg/router"
)

// HostMatcher is a middleware that matches any host with the given list of hosts.
// It also supports regex host like *.example.com
type HostMatcher struct {
	plainHosts    map[string]bool
	wildcardHosts []*regexp.Regexp
}

// NewHostMatcher creates a new instance of HostMatcher
func NewHostMatcher() *HostMatcher {
	return &HostMatcher{plainHosts: make(map[string]bool)}
}

// GetName retrieves the plugin's name
func (h *HostMatcher) GetName() string {
	return "host_matcher"
}

// GetMiddlewares retrieves the plugin's middlewares
func (h *HostMatcher) GetMiddlewares(config api.Config, referenceSpec *api.Spec) ([]router.Constructor, error) {
	var hosts []string
	aInterface := config["hosts"].([]interface{})
	for _, v := range aInterface {
		hosts = append(hosts, v.(string))
	}

	return []router.Constructor{
		middleware.NewHostMatcher(hosts).Handler,
	}, nil
}
