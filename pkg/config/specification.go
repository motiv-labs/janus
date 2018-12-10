package config

import (
	"time"

	"github.com/hellofresh/logging-go"
	"github.com/kelseyhightower/envconfig"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Specification for basic configurations
type Specification struct {
	Port                 int           `envconfig:"PORT"`
	GraceTimeOut         int64         `envconfig:"GRACE_TIMEOUT"`
	MaxIdleConnsPerHost  int           `envconfig:"MAX_IDLE_CONNS_PER_HOST"`
	BackendFlushInterval time.Duration `envconfig:"BACKEND_FLUSH_INTERVAL"`
	IdleConnTimeout      time.Duration `envconfig:"IDLE_CONN_TIMEOUT"`
	RequestID            bool          `envconfig:"REQUEST_ID_ENABLED"`
	Log                  logging.LogConfig
	Web                  Web
	Database             Database
	Stats                Stats
	Tracing              Tracing
	TLS                  TLS
	Cluster              Cluster
	RespondingTimeouts   RespondingTimeouts
}

// Cluster represents the cluster configuration
type Cluster struct {
	UpdateFrequency time.Duration `envconfig:"BACKEND_UPDATE_FREQUENCY"`
}

// RespondingTimeouts contains timeout configurations for incoming requests to the Janus instance.
type RespondingTimeouts struct {
	ReadTimeout  time.Duration `envconfig:"RESPONDING_TIMEOUTS_READ_TIMEOUT"`
	WriteTimeout time.Duration `envconfig:"RESPONDING_TIMEOUTS_WRITE_TIMEOUT"`
	IdleTimeout  time.Duration `envconfig:"RESPONDING_TIMEOUTS_IDLE_TIMEOUT"`
}

// Web represents the API configurations
type Web struct {
	Port        int `envconfig:"API_PORT"`
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

// Database holds the configuration for a database
type Database struct {
	DSN string `envconfig:"DATABASE_DSN"`
}

// Stats holds the configuration for stats
type Stats struct {
	DSN                   string   `envconfig:"STATS_DSN"`
	IDs                   string   `envconfig:"STATS_IDS"`
	AutoDiscoverThreshold uint     `envconfig:"STATS_AUTO_DISCOVER_THRESHOLD"`
	AutoDiscoverWhiteList []string `envconfig:"STATS_AUTO_DISCOVER_WHITE_LIST"`
	ErrorsSection         string   `envconfig:"STATS_ERRORS_SECTION"`
	Exporter              string   `envconfig:"STATS_EXPORTER"`
}

// Credentials represents the credentials that are going to be
// used by admin JWT configuration
type Credentials struct {
	// Algorithm defines admin JWT signing algorithm.
	// Currently the following algorithms are supported: HS256, HS384, HS512.
	Algorithm      string        `envconfig:"ALGORITHM"`
	Secret         string        `envconfig:"SECRET"`
	JanusAdminTeam string        `envconfig:"JANUS_ADMIN_TEAM"`
	Timeout        time.Duration `envconfig:"TOKEN_TIMEOUT"`
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

// Tracing represents the distributed tracing configuration
type Tracing struct {
	Exporter         string        `envconfig:"TRACING_EXPORTER"`
	ServiceName      string        `envconfig:"TRACING_SERVICE_NAME"`
	SamplingStrategy string        `envconfig:"TRACING_SAMPLING_STRATEGY"`
	SamplingParam    float64       `envconfig:"TRACING_SAMPLING_PARAM"`
	JaegerTracing    JaegerTracing `mapstructure:"jaeger"`
}

// JaegerTracing holds the Jaeger tracing configuration
type JaegerTracing struct {
	SamplingServerURL string `envconfig:"TRACING_JAEGER_SAMPLING_SERVER_URL"`
}

func init() {
	serviceName := "janus"

	viper.SetDefault("port", "8080")
	viper.SetDefault("tls.port", "8433")
	viper.SetDefault("tls.redirect", true)
	viper.SetDefault("backendFlushInterval", "20ms")
	viper.SetDefault("requestID", true)

	viper.SetDefault("respondingTimeouts.IdleTimeout", 180*time.Second)

	viper.SetDefault("cluster.updateFrequency", "10s")
	viper.SetDefault("database.dsn", "file:///etc/janus")

	viper.SetDefault("web.port", "8081")
	viper.SetDefault("web.tls.port", "8444")
	viper.SetDefault("web.tls.redirect", true)
	viper.SetDefault("web.credentials.algorithm", "HS256")
	viper.SetDefault("web.credentials.timeout", time.Hour)
	viper.SetDefault("web.credentials.basic.users", map[string]string{"admin": "admin"})
	viper.SetDefault("web.credentials.github.teams", make(map[string]string))

	viper.SetDefault("stats.dsn", "log://")
	viper.SetDefault("stats.errorsSection", "error-log")
	viper.SetDefault("stats.namespace", serviceName)

	viper.SetDefault("tracing.serviceName", serviceName)
	viper.SetDefault("tracing.samplingStrategy", "probabilistic")
	viper.SetDefault("tracing.samplingParam", 0.15)

	logging.InitDefaults(viper.GetViper(), "log")
}

//Load configuration variables
func Load(configFile string) (*Specification, error) {
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		dir, err := homedir.Dir()
		if err != nil {
			return nil, err
		}

		viper.SetConfigName("janus")
		viper.AddConfigPath(".")
		viper.AddConfigPath(dir)
		viper.AddConfigPath("/etc/janus")
	}

	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "No config file found")
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
	var err error

	// ensure the defaults are loaded
	err = viper.Unmarshal(&config)
	if err != nil {
		log.WithError(err).Warn("Failed unmarshaling config")
	}

	err = envconfig.Process("", &config)
	return &config, err
}
