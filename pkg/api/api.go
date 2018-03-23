package api

import (
	"encoding/json"

	"github.com/asaskevich/govalidator"
	"github.com/hellofresh/janus/pkg/proxy"
)

// Spec Holds an api definition and basic options
type Spec struct {
	*Definition
}

// Plugin represents the plugins for an API
type Plugin struct {
	Name    string                 `bson:"name" json:"name"`
	Enabled bool                   `bson:"enabled" json:"enabled"`
	Config  map[string]interface{} `bson:"config" json:"config"`
}

// Definition Represents an API that you want to proxy
type Definition struct {
	Name           string            `bson:"name" json:"name" valid:"required~name is required,matches(^[A-Za-z0-9]+(?:-[A-Za-z0-9]+)*$)~name cannot contain non-URL friendly characters"`
	Active         bool              `bson:"active" json:"active"`
	Proxy          *proxy.Definition `bson:"proxy" json:"proxy" valid:"required"`
	Plugins        []Plugin          `bson:"plugins" json:"plugins"`
	HealthCheck    HealthCheck       `bson:"health_check" json:"health_check"`
	CircuitBreaker CircuitBreaker    `bson:"circuit_breaker" json:"circuit_breaker"`
}

// CircuitBreaker is the configurations for the circuit breaker
type CircuitBreaker struct {
	Timeout                int `bson:"timeout" json:"timeout"`
	MaxConcurrentRequests  int `bson:"max_concurrent_requests" json:"max_concurrent_requests"`
	RequestVolumeThreshold int `bson:"request_volume_threshold" json:"request_volume_threshold"`
	SleepWindow            int `bson:"sleep_window" json:"sleep_window"`
	ErrorPercentThreshold  int `bson:"error_percent_threshold" json:"error_percent_threshold"`
}

// HealthCheck represents the health check configs
type HealthCheck struct {
	URL     string `bson:"url" json:"url" valid:"url"`
	Timeout int    `bson:"timeout" json:"timeout"`
}

// NewDefinition creates a new API Definition with default values
func NewDefinition() *Definition {
	return &Definition{
		Active:  true,
		Plugins: make([]Plugin, 0),
		Proxy:   proxy.NewDefinition(),
	}
}

// Validate validates proxy data
func (d *Definition) Validate() (bool, error) {
	return govalidator.ValidateStruct(d)
}

// UnmarshalJSON api.Definition JSON.Unmarshaller implementation
func (d *Definition) UnmarshalJSON(b []byte) error {
	// Aliasing Definition to avoid recursive call of this method
	type definitionAlias Definition
	defAlias := definitionAlias(*NewDefinition())

	if err := json.Unmarshal(b, &defAlias); err != nil {
		return err
	}

	*d = Definition(defAlias)
	return nil
}
