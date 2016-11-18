package config

import "github.com/kelseyhightower/envconfig"

// Specification for basic configurations
type Specification struct {
	DatabaseDSN string `envconfig:"DATABASE_DSN"`
	Port        int    `envconfig:"PORT"`
	Debug       bool   `envconfig:"DEBUG"`
	StorageDSN  string `envconfig:"REDIS_DSN"`
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
