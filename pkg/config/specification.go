package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Specification for basic configurations
type Specification struct {
	Port                 int           `envconfig:"PORT" default:"8080" description:"Default application port"`
	APIPort              int           `envconfig:"API_PORT" default:"8081" description:"Admin API port"`
	Debug                bool          `envconfig:"DEBUG" description:"Enable debug mode"`
	LogLevel             string        `envconfig:"LOG_LEVEL" default:"info" description:"Log level"`
	GraceTimeOut         int64         `envconfig:"GRACE_TIMEOUT" description:"Duration to give active requests a chance to finish during hot-reload"`
	MaxIdleConnsPerHost  int           `envconfig:"MAX_IDLE_CONNS_PER_HOST" description:"If non-zero, controls the maximum idle (keep-alive) to keep per-host."`
	InsecureSkipVerify   bool          `envconfig:"INSECURE_SKIP_VERIFY" description:"Disable SSL certificate verification"`
	StorageDSN           string        `envconfig:"STORAGE_DSN" default:"memory://localhost" description:"The Storage DSN, this could be 'memory' or 'redis'"`
	CertPathTLS          string        `envconfig:"CERT_PATH" description:"Path of certificate when using TLS"`
	KeyPathTLS           string        `envconfig:"KEY_PATH" description:"Path of key when using TLS"`
	BackendFlushInterval time.Duration `envconfig:"BACKEND_FLUSH_INTERVAL" default:"20ms" description:"Flush interval for upgraded Proxy connections"`
	CloseIdleConnsPeriod time.Duration `envconfig:"CLOSE_IDLE_CONNS_PERIOD" description:"Defines the time period of how often the idle connections maintained by the proxy are closed."`
	Database             Database
	Statsd               Statsd
	Credentials          Credentials
	Tracing              Tracing
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

// Credentials represents the credentials that are going to be
// used by JWT configuration
type Credentials struct {
	Secret   string `envconfig:"SECRET" required:"true"`
	Username string `envconfig:"ADMIN_USERNAME" default:"admin"`
	Password string `envconfig:"ADMIN_PASSWORD" default:"admin"`
}

// GoogleCloudTracing holds the Google Application Default Credentials
type GoogleCloudTracing struct {
	ProjectID    string `envconfig:"TRACING_GC_PROJECT_ID"`
	Email        string `envconfig:"TRACING_GC_EMAIL"`
	PrivateKey   string `envconfig:"TRACING_GC_PRIVATE_KEY"`
	PrivateKeyID string `envconfig:"TRACING_GC_PRIVATE_ID"`
}

// AppdashTracing holds the Appdash tracing configuration
type AppdashTracing struct {
	DSN string `envconfig:"TRACING_APPDASH_DSN"`
	URL string `envconfig:"TRACING_APPDASH_URL"`
}

// Tracing represents the distributed tracing configuration
type Tracing struct {
	GoogleCloudTracing GoogleCloudTracing
	AppdashTracing     AppdashTracing
}

// IsGoogleCloudEnabled checks if google cloud is enabled
func (t Tracing) IsGoogleCloudEnabled() bool {
	return len(t.GoogleCloudTracing.Email) > 0 && len(t.GoogleCloudTracing.PrivateKey) > 0 && len(t.GoogleCloudTracing.PrivateKeyID) > 0
}

// IsAppdashEnabled checks if appdash is enabled
func (t Tracing) IsAppdashEnabled() bool {
	return len(t.AppdashTracing.DSN) > 0
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
