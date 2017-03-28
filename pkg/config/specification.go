package config

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/fsnotify/fsnotify"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/viper"
)

// Specification for basic configurations
type Specification struct {
	Port                 int           `envconfig:"PORT" description:"Default application port"`
	APIPort              int           `envconfig:"API_PORT" mapstructure:"api_port" description:"Admin API port"`
	Debug                bool          `envconfig:"DEBUG" description:"Enable debug mode"`
	LogLevel             string        `envconfig:"LOG_LEVEL" mapstructure:"log_level" description:"Log level"`
	GraceTimeOut         int64         `envconfig:"GRACE_TIMEOUT" mapstructure:"grace_timeout" description:"Duration to give active requests a chance to finish during hot-reload"`
	MaxIdleConnsPerHost  int           `envconfig:"MAX_IDLE_CONNS_PER_HOST" mapstructure:"max_idle_conns_per_host" description:"If non-zero, controls the maximum idle (keep-alive) to keep per-host."`
	InsecureSkipVerify   bool          `envconfig:"INSECURE_SKIP_VERIFY" mapstructure:"insecure_skip_verify" description:"Disable SSL certificate verification"`
	CertPathTLS          string        `envconfig:"CERT_PATH" mapstructure:"cert_path" description:"Path of certificate when using TLS"`
	KeyPathTLS           string        `envconfig:"KEY_PATH" mapstructure:"key_path" description:"Path of key when using TLS"`
	BackendFlushInterval time.Duration `envconfig:"BACKEND_FLUSH_INTERVAL" mapstructure:"backend_flush_interval" description:"Flush interval for upgraded Proxy connections"`
	CloseIdleConnsPeriod time.Duration `envconfig:"CLOSE_IDLE_CONNS_PERIOD" mapstructure:"close_idle_conns_period" description:"Defines the time period of how often the idle connections maintained by the proxy are closed."`
	Database             Database
	Storage              Storage
	Stats                Stats
	Credentials          Credentials
	Tracing              Tracing
}

// IsHTTPS checks if you have https enabled
func (s *Specification) IsHTTPS() bool {
	return s.CertPathTLS != "" && s.KeyPathTLS != ""
}

// Storage holds the configuration for a storage
type Storage struct {
	DSN string `envconfig:"STORAGE_DSN" description:"The Storage DSN, this could be 'memory' or 'redis'"`
}

// Database holds the configuration for a database
type Database struct {
	DSN string `envconfig:"DATABASE_DSN"`
}

// Stats holds the configuration for stats
type Stats struct {
	DSN    string `envconfig:"STATS_DSN"`
	Prefix string `envconfig:"STATS_PREFIX"`
	IDs    string `envconfig:"STATS_IDS"`
}

// Credentials represents the credentials that are going to be
// used by JWT configuration
type Credentials struct {
	Secret   string `envconfig:"SECRET"`
	Username string `envconfig:"ADMIN_USERNAME"`
	Password string `envconfig:"ADMIN_PASSWORD"`
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
	GoogleCloudTracing GoogleCloudTracing `mapstructure:"google_cloud"`
	AppdashTracing     AppdashTracing     `mapstructure:"appdash"`
}

// IsGoogleCloudEnabled checks if google cloud is enabled
func (t Tracing) IsGoogleCloudEnabled() bool {
	return len(t.GoogleCloudTracing.Email) > 0 && len(t.GoogleCloudTracing.PrivateKey) > 0 && len(t.GoogleCloudTracing.PrivateKeyID) > 0 && len(t.GoogleCloudTracing.ProjectID) > 0
}

// IsAppdashEnabled checks if appdash is enabled
func (t Tracing) IsAppdashEnabled() bool {
	return len(t.AppdashTracing.DSN) > 0
}

func init() {
	viper.SetDefault("port", "8080")
	viper.SetDefault("api_port", "8081")
	viper.SetDefault("log_level", "info")
	viper.SetDefault("backend_flush_interval", "20ms")
	viper.SetDefault("database.dsn", "file:///etc/janus")
	viper.SetDefault("storage.dsn", "memory://localhost")
	viper.SetDefault("credentials.username", "admin")
	viper.SetDefault("credentials.password", "admin")
}

//Load configuration variables
func Load() (*Specification, error) {
	var config Specification

	// from a config file
	viper.SetConfigName("janus")
	viper.AddConfigPath("/etc/janus")
	viper.AddConfigPath(".")

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Debug("Configuration changed")
	})

	if err := viper.ReadInConfig(); err != nil {
		log.WithError(err).Warn("No config file found")
		return LoadEnv()
	}

	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

//LoadEnv loads configuration from environment variables
func LoadEnv() (*Specification, error) {
	var config Specification

	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	err := envconfig.Process("", &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
