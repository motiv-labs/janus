package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Specification for basic configurations
type Specification struct {
	Port                int    `envconfig:"PORT" default:"8080"`
	APIPort             int    `envconfig:"API_PORT" default:"8081"`
	Debug               bool   `envconfig:"DEBUG"`
	LogLevel            string `envconfig:"LOG_LEVEL" default:"info"`
	GraceTimeOut        int64  `envconfig:"GRACE_TIMEOUT"`
	MaxIdleConnsPerHost int    `envconfig:"MAX_IDLE_CONNS_PER_HOST"`
	InsecureSkipVerify  bool   `envconfig:"INSECURE_SKIP_VERIFY"`
	// The Storage DSN, this could be `memory` or `redis`
	StorageDSN string `envconfig:"STORAGE_DSN" default:"memory://localhost"`

	// Path of certificate when using TLS
	CertPathTLS string `envconfig:"CERT_PATH"`

	// Path of key when using TLS
	KeyPathTLS string `envconfig:"KEY_PATH"`

	// Flush interval for upgraded Proxy connections
	BackendFlushInterval time.Duration `envconfig:"BACKEND_FLUSH_INTERVAL" default:"20ms"`

	// Defines the time period of how often the idle connections maintained
	// by the proxy are closed.
	CloseIdleConnsPeriod time.Duration `envconfig:"CLOSE_IDLE_CONNS_PERIOD"`

	Database    Database
	Statsd      Statsd
	Credentials Credentials
	Application Application
}

// IsHTTPS checks if you have https enabled
func (s *Specification) IsHTTPS() bool {
	return s.CertPathTLS != "" && s.KeyPathTLS != ""
}

// Database holds the configuration for a database
type Database struct {
	DSN string `envconfig:"DATABASE_DSN" default:"file:///etc/janus"`
}

// Statsd holds the configuration for statsd
type Statsd struct {
	DSN    string `envconfig:"STATSD_DSN"`
	Prefix string `envconfig:"STATSD_PREFIX"`
}

// IsEnabled checks if you have metrics enabled
func (s Statsd) IsEnabled() bool {
	return len(s.DSN) == 0
}

// HasPrefix checks if you have any prefix configured
func (s Statsd) HasPrefix() bool {
	return len(s.Prefix) > 0
}

// Application represents a simple application definition
type Application struct {
	Name string `envconfig:"APP_NAME" default:"Janus"`
}

// Credentials represents the credentials that are going to be
// used by JWT configuration
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
