package proxy

import (
	"encoding/json"

	"github.com/asaskevich/govalidator"
	"github.com/hellofresh/janus/pkg/router"
)

// Definition defines proxy rules for a route
type Definition struct {
	PreserveHost        bool     `bson:"preserve_host" json:"preserve_host" mapstructure:"preserve_host"`
	ListenPath          string   `bson:"listen_path" json:"listen_path" mapstructure:"listen_path" valid:"required"`
	UpstreamURL         string   `bson:"upstream_url" json:"upstream_url" mapstructure:"upstream_url" valid:"url,required"`
	InsecureSkipVerify  bool     `bson:"insecure_skip_verify" json:"insecure_skip_verify" mapstructure:"insecure_skip_verify"`
	StripPath           bool     `bson:"strip_path" json:"strip_path" mapstructure:"strip_path"`
	AppendPath          bool     `bson:"append_path" json:"append_path" mapstructure:"append_path"`
	EnableLoadBalancing bool     `bson:"enable_load_balancing" json:"enable_load_balancing" mapstructure:"enable_load_balancing"`
	Methods             []string `bson:"methods" json:"methods"`
	Hosts               []string `bson:"hosts" json:"hosts"`
}

// NewDefinition creates a new Proxy Definition with default values
func NewDefinition() *Definition {
	return &Definition{
		Methods: make([]string, 0),
		Hosts:   make([]string, 0),
	}
}

// Validate validates proxy data
func (d *Definition) Validate() (bool, error) {
	return govalidator.ValidateStruct(d)
}

// Route is the container for a proxy and it's handlers
type Route struct {
	Proxy    *Definition
	Inbound  InChain
	Outbound OutChain
}

type routeJSONProxy struct {
	Proxy *Definition `json:"proxy"`
}

// NewRoute creates an instance of Route
func NewRoute(proxy *Definition) *Route {
	return &Route{Proxy: proxy}
}

// NewRouteWithInOut creates an instance of Route with inbound and outbound handlers
func NewRouteWithInOut(proxy *Definition, inbound InChain, outbound OutChain) *Route {
	return &Route{proxy, inbound, outbound}
}

// AddInbound adds inbound middlewares
func (r *Route) AddInbound(in ...router.Constructor) {
	for _, i := range in {
		r.Inbound = append(r.Inbound, i)
	}
}

// AddOutbound adds outbound middlewares
func (r *Route) AddOutbound(out ...OutLink) {
	for _, o := range out {
		r.Outbound = append(r.Outbound, o)
	}
}

// JSONMarshal encodes route struct to JSON
func (r *Route) JSONMarshal() ([]byte, error) {
	return json.Marshal(routeJSONProxy{r.Proxy})
}

// JSONUnmarshalRoute decodes route struct from JSON
func JSONUnmarshalRoute(rawRoute []byte) (*Route, error) {
	var proxyRoute routeJSONProxy
	if err := json.Unmarshal(rawRoute, &proxyRoute); err != nil {
		return nil, err
	}
	return NewRoute(proxyRoute.Proxy), nil
}
