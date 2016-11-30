package config

import "github.com/kelseyhightower/envconfig"

// Specification for basic configurations
type Specification struct {
	DatabaseDSN  string `envconfig:"DATABASE_DSN"`
	Port         int    `envconfig:"PORT"`
	Debug        bool   `envconfig:"DEBUG"`
	StatsdDSN    string `envconfig:"STATSD_DSN"`
	StatsdPrefix string `envconfig:"STATSD_PREFIX"`
	StorageDSN   string `envconfig:"REDIS_DSN"`
	Credentials  Credentials
	Application  Application
}

type Application struct {
	Name    string `envconfig:"APP_NAME" default:"Janus"`
	Version string `envconfig:"APP_VERSION" default:"1.0"`
}

type Credentials struct {
	Secret   string `envconfig:"SECRET"`
	Username string `envconfig:"ADMIN_USERNAME"`
	Password string `envconfig:"ADMIN_PASSWORD"`
}

//LoadEnv loads environment variables
func LoadEnv() (*Specification, error) {
	var config Specification
	err := envconfig.Process("", &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
