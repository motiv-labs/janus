package config

import (
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/hellofresh/logging-go"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Specification for basic configurations
type Specification struct {
	Port                 int           `envconfig:"PORT"`
	Debug                bool          `envconfig:"DEBUG"`
	GraceTimeOut         int64         `envconfig:"GRACE_TIMEOUT"`
	MaxIdleConnsPerHost  int           `envconfig:"MAX_IDLE_CONNS_PER_HOST"`
	InsecureSkipVerify   bool          `envconfig:"INSECURE_SKIP_VERIFY"`
	BackendFlushInterval time.Duration `envconfig:"BACKEND_FLUSH_INTERVAL"`
	CloseIdleConnsPeriod time.Duration `envconfig:"CLOSE_IDLE_CONNS_PERIOD"`
	CertFile             string        `envconfig:"CERT_PATH"`
	KeyFile              string        `envconfig:"KEY_PATH"`
	Log                  logging.LogConfig
	Web                  Web
	Database             Database
	Storage              Storage
	Stats                Stats
	Tracing              Tracing
}

// Web represents the API configurations
type Web struct {
	Port        int    `envconfig:"API_PORT"`
	CertFile    string `envconfig:"API_CERT_PATH"`
	KeyFile     string `envconfig:"API_KEY_PATH"`
	ReadOnly    bool   `envconfig:"API_READONLY"`
	Credentials Credentials
}

// IsHTTPS checks if you have https enabled
func (s *Web) IsHTTPS() bool {
	return s.CertFile != "" && s.KeyFile != ""
}

// Storage holds the configuration for a storage
type Storage struct {
	DSN string `envconfig:"STORAGE_DSN"`
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
	GoogleCloudTracing GoogleCloudTracing `mapstructure:"googleCloud"`
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
	viper.SetDefault("backendFlushInterval", "20ms")
	viper.SetDefault("database.dsn", "file:///etc/janus")
	viper.SetDefault("storage.dsn", "memory://localhost")
	viper.SetDefault("web.port", "8081")
	viper.SetDefault("web.credentials.username", "admin")
	viper.SetDefault("web.credentials.password", "admin")

	logging.InitDefaults(viper.GetViper(), "log")
}

//Load configuration variables
func Load(configFile string) (*Specification, error) {
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.SetConfigName("janus")
		viper.AddConfigPath("/etc/janus")
		viper.AddConfigPath(".")
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Debug("Configuration changed")
	})

	if err := viper.ReadInConfig(); err != nil {
		log.WithError(err).Warn("No config file found")
		return LoadEnv()
	}

	var config Specification
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
