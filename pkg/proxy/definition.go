package proxy

import (
	"encoding/json"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/hellofresh/janus/pkg/router"
)

// Route is the container for a proxy and it's handlers
type Route struct {
	proxy    *Definition
	handlers []router.Constructor
}

type routeJSONProxy struct {
	Proxy *Definition `json:"proxy"`
}

// NewRoute creates an instance of Route
func NewRoute(proxy *Definition, handlers ...router.Constructor) *Route {
	return &Route{proxy, handlers}
}

// JSONMarshal encodes route struct to JSON
func (r *Route) JSONMarshal() ([]byte, error) {
	return json.Marshal(routeJSONProxy{r.proxy})
}

// JSONUnmarshalRoute decodes route struct from JSON
func JSONUnmarshalRoute(rawRoute []byte) (*Route, error) {
	var proxyRoute routeJSONProxy
	if err := json.Unmarshal(rawRoute, &proxyRoute); err != nil {
		return nil, err
	}
	return NewRoute(proxyRoute.Proxy), nil
}

// Definition defines proxy rules for a route
type Definition struct {
	PreserveHost        bool     `bson:"preserve_host" json:"preserve_host"`
	ListenPath          string   `bson:"listen_path" json:"listen_path" valid:"required"`
	UpstreamURL         string   `bson:"upstream_url" json:"upstream_url" valid:"url,required"`
	StripPath           bool     `bson:"strip_path" json:"strip_path"`
	AppendPath          bool     `bson:"append_path" json:"append_path"`
	EnableLoadBalancing bool     `bson:"enable_load_balancing" json:"enable_load_balancing"`
	Methods             []string `bson:"methods" json:"methods"`
	Hosts               []string `bson:"hosts" json:"hosts"`
}

// Validate validates proxy data
func Validate(proxy *Definition) bool {
	if nil == proxy {
		return false
	}

	if proxy.ListenPath == "" {
		log.Warning("Listen path is empty")
		return false
	}

	if strings.Contains(proxy.ListenPath, " ") {
		log.Warning("Listen path contains spaces, is invalid")
		return false
	}

	return true
}
