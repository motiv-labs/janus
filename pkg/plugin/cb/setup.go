package cb

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/afex/hystrix-go/hystrix"
	metricCollector "github.com/afex/hystrix-go/hystrix/metric_collector"
	"github.com/afex/hystrix-go/plugins"
	"github.com/asaskevich/govalidator"
	log "github.com/sirupsen/logrus"

	"github.com/hellofresh/janus/pkg/plugin"
	"github.com/hellofresh/janus/pkg/proxy"
)

const (
	pluginName = "cb"
)

// Config represents the Body Limit configuration
type Config struct {
	hystrix.CommandConfig
	Name      string `json:"name"`
	Predicate string `json:"predicate"`
}

func init() {
	plugin.RegisterEventHook(plugin.StartupEvent, onStartup)
	plugin.RegisterEventHook(plugin.AdminAPIStartupEvent, onAdminAPIStartup)
	plugin.RegisterPlugin(pluginName, plugin.Plugin{
		Action:   setupCB,
		Validate: validateConfig,
	})
}

func setupCB(def *proxy.RouterDefinition, rawConfig plugin.Config) error {
	var c Config
	err := plugin.Decode(rawConfig, &c)
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"plugin_event": plugin.SetupEvent,
		"plugin":       pluginName,
		"name":         c.Name,
	}).Debug("Configuring cb plugin")

	hystrix.ConfigureCommand(c.Name, hystrix.CommandConfig{
		Timeout:               c.Timeout,
		MaxConcurrentRequests: c.MaxConcurrentRequests,
		ErrorPercentThreshold: c.ErrorPercentThreshold,
		SleepWindow:           c.SleepWindow,
	})

	def.AddMiddleware(NewCBMiddleware(c))
	return nil
}

func validateConfig(rawConfig plugin.Config) (bool, error) {
	var config Config
	err := plugin.Decode(rawConfig, &config)
	if err != nil {
		return false, err
	}

	return govalidator.ValidateStruct(config)
}

func onAdminAPIStartup(event interface{}) error {
	logger := log.WithFields(log.Fields{
		"plugin_event": plugin.AdminAPIStartupEvent,
		"plugin":       pluginName,
	})

	e, ok := event.(plugin.OnAdminAPIStartup)
	if !ok {
		return errors.New("could not convert event to admin startup type")
	}

	logger.Debug("Registering hystrix stream endpoint")
	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()

	e.Router.GET("/hystrix", hystrixStreamHandler.ServeHTTP)
	return nil
}

func onStartup(event interface{}) error {
	logger := log.WithFields(log.Fields{
		"plugin_event": plugin.StartupEvent,
		"plugin":       pluginName,
	})

	e, ok := event.(plugin.OnStartup)
	if !ok {
		return errors.New("could not convert event to startup type")
	}

	logger.WithField("metrics_dsn", e.Config.Stats.DSN).Debug("Statsd metrics enabled")

	dsnURL, err := url.Parse(e.Config.Stats.DSN)
	if err != nil {
		return fmt.Errorf("could not parse stats dsn: %w", err)
	}

	c, err := plugins.InitializeStatsdCollector(&plugins.StatsdCollectorConfig{
		StatsdAddr: dsnURL.Host,
		Prefix:     strings.Trim(dsnURL.Path, "/"),
	})
	if err != nil {
		return fmt.Errorf("could not initialize statsd client: %w", err)
	}

	metricCollector.Registry.Register(c.NewStatsdCollector)
	logger.Debug("Metrics enabled")

	return nil
}
