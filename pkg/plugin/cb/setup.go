package cb

import (
	"github.com/afex/hystrix-go/hystrix"
	"github.com/afex/hystrix-go/hystrix/metric_collector"
	"github.com/afex/hystrix-go/plugins"
	"github.com/hellofresh/janus/pkg/api"
	"github.com/hellofresh/janus/pkg/plugin"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const (
	statsdPrefix = "hystrix"
	pluginName   = "cb"
)

// Config represents the Body Limit configuration
type Config struct {
	hystrix.CommandConfig
	Predicate string `json:"predicate"`
}

func init() {
	plugin.RegisterEventHook(plugin.StartupEvent, onStartup)
	plugin.RegisterEventHook(plugin.AdminAPIStartupEvent, onAdminAPIStartup)
	plugin.RegisterPlugin(pluginName, plugin.Plugin{
		Action: setupCB,
	})
}

func setupCB(def *api.Definition, rawConfig plugin.Config) error {
	logger := log.WithFields(log.Fields{
		"plugin_event": plugin.SetupEvent,
		"plugin":       pluginName,
	})

	var c Config
	err := plugin.Decode(rawConfig, &c)
	if err != nil {
		return err
	}

	logger.WithField("name", def.Name).Debug("Configuring cb plugin")
	hystrix.ConfigureCommand(def.Name, hystrix.CommandConfig{
		Timeout:               c.Timeout,
		MaxConcurrentRequests: c.MaxConcurrentRequests,
		ErrorPercentThreshold: c.ErrorPercentThreshold,
		SleepWindow:           c.SleepWindow,
	})

	def.Proxy.AddMiddleware(NewCBMiddleware(c, def))
	return nil
}

func onAdminAPIStartup(event interface{}) error {
	logger := log.WithFields(log.Fields{
		"plugin_event": plugin.AdminAPIStartupEvent,
		"plugin":       pluginName,
	})

	e, ok := event.(plugin.OnAdminAPIStartup)
	if !ok {
		return errors.New("Could not convert event to admin startup type")
	}

	logger.Debug("Registring hystrix stream endpoint")
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
		return errors.New("Could not convert event to startup type")
	}

	logger.WithField("metrics_dsn", e.Config.Stats.DSN).Debug("Statsd metrics enabled")
	c, err := plugins.InitializeStatsdCollector(&plugins.StatsdCollectorConfig{
		StatsdAddr: e.Config.Stats.DSN,
		Prefix:     statsdPrefix,
	})
	if err != nil {
		return errors.Wrap(err, "could not initialize statsd client")
	}

	metricCollector.Registry.Register(c.NewStatsdCollector)
	logger.Debug("Metrics enabled")

	return nil
}
