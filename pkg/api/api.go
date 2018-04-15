package api

import (
	"encoding/json"

	"github.com/asaskevich/govalidator"
	"github.com/hellofresh/janus/pkg/proxy"
)

// Plugin represents the plugins for an API
type Plugin struct {
	Name    string                 `bson:"name" json:"name"`
	Enabled bool                   `bson:"enabled" json:"enabled"`
	Config  map[string]interface{} `bson:"config" json:"config"`
}

// Definition Represents an API that you want to proxy
type Definition struct {
	Name        string            `bson:"name" json:"name" valid:"required~name is required,matches(^[A-Za-z0-9]+(?:-[A-Za-z0-9]+)*$)~name cannot contain non-URL friendly characters"`
	Active      bool              `bson:"active" json:"active"`
	Proxy       *proxy.Definition `bson:"proxy" json:"proxy" valid:"required"`
	Plugins     []Plugin          `bson:"plugins" json:"plugins"`
	HealthCheck HealthCheck       `bson:"health_check" json:"health_check"`
}

// HealthCheck represents the health check configs
type HealthCheck struct {
	URL     string `bson:"url" json:"url" valid:"url"`
	Timeout int    `bson:"timeout" json:"timeout"`
}

// Configuration represents all the api definitions
type Configuration struct {
	Definitions []*Definition
}

// ConfigurationChanged is the message that is sent when a database configuration has changed
type ConfigurationChanged struct {
	Configurations *Configuration
}

// ConfigurationOperation is the available operations that a configuration can have
type ConfigurationOperation int

// ConfigurationMessage is used to notify listeners about something that happened with a configuration
type ConfigurationMessage struct {
	Operation     ConfigurationOperation
	Configuration *Definition
}

const (
	// RemovedOperation means a definition was removed
	RemovedOperation ConfigurationOperation = iota
	// UpdatedOperation means a definition was updated
	UpdatedOperation
	// AddedOperation means a definition was added
	AddedOperation
)

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
