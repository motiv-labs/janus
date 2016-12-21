package config

import "github.com/kelseyhightower/envconfig"

// Specification for basic configurations
type Specification struct {
	DatabaseDSN         string `envconfig:"DATABASE_DSN"`
	Port                int    `envconfig:"PORT"`
	Debug               bool   `envconfig:"DEBUG"`
	LogLevel            string `envconfig:"LOG_LEVEL" default:"info"`
	GraceTimeOut        int64  `envconfig:"GRACE_TIMEOUT"`
	MaxIdleConnsPerHost int    `description:"Disable SSL certificate verification" envconfig:"MAX_IDLE_CONNS_PER_HOST"`
	InsecureSkipVerify  bool   `description:"Disable SSL certificate verification"`
	StorageDSN          string `envconfig:"REDIS_DSN"`
	Statsd              Statsd
	Credentials         Credentials
	Application         Application
}

// Statsd holds the configuration for statsd
type Statsd struct {
	DSN    string `envconfig:"STATSD_DSN"`
	Prefix string `envconfig:"STATSD_PREFIX"`
}

type Application struct {
	Name    string `envconfig:"APP_NAME" default:"Janus"`
	Version string `envconfig:"APP_VERSION" default:"1.0"`
}

type Credentials struct {
	Secret   string `envconfig:"SECRET" required:"true"`
	Username string `envconfig:"ADMIN_USERNAME" default:"admin"`
	Password string `envconfig:"ADMIN_PASSWORD" default:"admin"`
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
