package config

import (
	"os/user"
	"time"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/hellofresh/logging-go"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// Specification for basic configurations
type Specification struct {
	Port                 int           `envconfig:"PORT"`
	Debug                bool          `envconfig:"DEBUG"`
	GraceTimeOut         int64         `envconfig:"GRACE_TIMEOUT"`
	MaxIdleConnsPerHost  int           `envconfig:"MAX_IDLE_CONNS_PER_HOST"`
	BackendFlushInterval time.Duration `envconfig:"BACKEND_FLUSH_INTERVAL"`
	CloseIdleConnsPeriod time.Duration `envconfig:"CLOSE_IDLE_CONNS_PERIOD"`
	Log                  logging.LogConfig
	Web                  Web
	Database             Database
	Storage              Storage
	Stats                Stats
	Tracing              Tracing
	TLS                  TLS
	CircuitBreaker       CircuitBreaker
}

// Web represents the API configurations
type Web struct {
	Port        int  `envconfig:"API_PORT"`
	ReadOnly    bool `envconfig:"API_READONLY"`
	Credentials Credentials
	TLS         TLS
}

// TLS represents the TLS configurations
type TLS struct {
	Port     int    `envconfig:"PORT"`
	CertFile string `envconfig:"CERT_PATH"`
	KeyFile  string `envconfig:"KEY_PATH"`
	Redirect bool   `envconfig:"REDIRECT"`
}

// IsHTTPS checks if you have https enabled
func (s *TLS) IsHTTPS() bool {
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
	DSN                   string   `envconfig:"STATS_DSN"`
	Prefix                string   `envconfig:"STATS_PREFIX"`
	IDs                   string   `envconfig:"STATS_IDS"`
	AutoDiscoverThreshold uint     `envconfig:"STATS_AUTO_DISCOVER_THRESHOLD"`
	AutoDiscoverWhiteList []string `envconfig:"STATS_AUTO_DISCOVER_WHITE_LIST"`
	ErrorsSection         string   `envconfig:"STATS_ERRORS_SECTION"`
}

// Credentials represents the credentials that are going to be
// used by admin JWT configuration
type Credentials struct {
	// Algorithm defines admin JWT signing algorithm.
	// Currently the following algorithms are supported: HS256, HS384, HS512.
	Algorithm      string `envconfig:"ALGORITHM"`
	Secret         string `envconfig:"SECRET"`
	JanusAdminTeam string `envconfig:"JANUS_ADMIN_TEAM"`
	Github         Github
	Basic          Basic
}

// Basic holds the basic users configurations
type Basic struct {
	Users map[string]string `envconfig:"BASIC_USERS"`
}

// Github holds the github configurations
type Github struct {
	Organizations []string          `envconfig:"GITHUB_ORGANIZATIONS"`
	Teams         map[string]string `envconfig:"GITHUB_TEAMS"`
}

// IsConfigured checks if github is enabled
func (auth *Github) IsConfigured() bool {
	return len(auth.Organizations) > 0 ||
		len(auth.Teams) > 0
}

// GoogleCloudTracing holds the Google Application Default Credentials
type GoogleCloudTracing struct {
	ProjectID    string `envconfig:"TRACING_GC_PROJECT_ID"`
	Email        string `envconfig:"TRACING_GC_EMAIL"`
	PrivateKey   string `envconfig:"TRACING_GC_PRIVATE_KEY"`
	PrivateKeyID string `envconfig:"TRACING_GC_PRIVATE_ID"`
}

// JaegerTracing holds the Jaeger tracing configuration
type JaegerTracing struct {
	DSN                 string `envconfig:"TRACING_JAEGER_DSN"`
	ServiceName         string `envconfig:"TRACING_JAEGER_SERVICE_NAME"`
	BufferFlushInterval string `envconfig:"TRACING_JAEGER_BUFFER_FLUSH_INTERVAL"`
	LogSpans            bool   `envconfig:"TRACING_JAEGER_LOG_SPANS"`
	QueueSize           int    `envconfig:"TRACING_JAEGER_QUEUE_SIZE"`
}

// Tracing represents the distributed tracing configuration
type Tracing struct {
	Provider           string             `envconfig:"TRACING_PROVIDER"`
	GoogleCloudTracing GoogleCloudTracing `mapstructure:"googleCloud"`
	JaegerTracing      JaegerTracing      `mapstructure:"jaeger"`
}

// CircuitBreaker represents the global circuit breaker settings
type CircuitBreaker struct {
	Timeout               int `envconfig:"CB_TIMEOUT"`
	MaxConcurrent         int `envconfig:"CB_MAX_CONCURRENT"`
	VolumeThreshold       int `envconfig:"CB_VOLUME_THRESHOLD"`
	SleepWindow           int `envconfig:"CB_SLEEP_WINDOW"`
	ErrorPercentThreshold int `envconfig:"CB_ERROR_PRECENT_THRESHOLD"`
}

func init() {
	serviceName := "janus"

	viper.SetDefault("port", "8080")
	viper.SetDefault("tls.port", "8433")
	viper.SetDefault("tls.redirect", true)
	viper.SetDefault("backendFlushInterval", "20ms")
	viper.SetDefault("database.dsn", "file:///etc/janus")
	viper.SetDefault("storage.dsn", "memory://localhost")
	viper.SetDefault("web.port", "8081")
	viper.SetDefault("web.tls.port", "8444")
	viper.SetDefault("web.tls.redisrect", true)
	viper.SetDefault("web.credentials.algorithm", "HS256")
	viper.SetDefault("web.credentials.basic.users", map[string]string{"admin": "admin"})
	viper.SetDefault("stats.dsn", "log://")
	viper.SetDefault("stats.errorsSection", "error-log")
	viper.SetDefault("tracing.jaeger.serviceName", serviceName)
	viper.SetDefault("tracing.jaeger.bufferFlushInterval", "1s")
	viper.SetDefault("tracing.jaeger.logSpans", false)
	viper.SetDefault("circuitBreaker.timeout", hystrix.DefaultTimeout)
	viper.SetDefault("circuitBreaker.MaxConcurrent", hystrix.DefaultMaxConcurrent)
	viper.SetDefault("circuitBreaker.VolumeThreshold", hystrix.DefaultVolumeThreshold)
	viper.SetDefault("circuitBreaker.SleepWindow", hystrix.DefaultSleepWindow)
	viper.SetDefault("circuitBreaker.ErrorPercentThreshold", hystrix.DefaultErrorPercentThreshold)

	logging.InitDefaults(viper.GetViper(), "log")
}

//Load configuration variables
func Load(configFile string) (*Specification, error) {
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		usr, err := user.Current()
		if err != nil {
			return nil, err
		}

		viper.SetConfigName("janus")
		viper.AddConfigPath(usr.HomeDir)
		viper.AddConfigPath("/etc/janus")
		viper.AddConfigPath(".")
	}

	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "No config file found")
	}

	var config Specification
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	hystrix.DefaultTimeout = config.CircuitBreaker.Timeout
	hystrix.DefaultMaxConcurrent = config.CircuitBreaker.MaxConcurrent
	hystrix.DefaultVolumeThreshold = config.CircuitBreaker.VolumeThreshold
	hystrix.DefaultSleepWindow = config.CircuitBreaker.SleepWindow
	hystrix.DefaultErrorPercentThreshold = config.CircuitBreaker.ErrorPercentThreshold

	return &config, nil
}

//LoadEnv loads configuration from environment variables
func LoadEnv() (*Specification, error) {
	var config Specification

	// ensure the defaults are loaded
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	err := envconfig.Process("", &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
