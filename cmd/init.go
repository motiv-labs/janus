package cmd

import (
	"os"
	"path/filepath"
	"time"

	"github.com/hellofresh/janus/pkg/config"
	obs "github.com/hellofresh/janus/pkg/observability"
	"github.com/hellofresh/stats-go"
	"github.com/hellofresh/stats-go/bucket"
	"github.com/hellofresh/stats-go/client"
	"github.com/hellofresh/stats-go/hooks"
	log "github.com/sirupsen/logrus"
	"go.opencensus.io/exporter/jaeger"
	"go.opencensus.io/exporter/prometheus"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
)

var (
	globalConfig *config.Specification
	statsClient  client.Client
)

func initConfig() {
	var err error
	globalConfig, err = config.Load(configFile)
	if nil != err {
		log.WithError(err).Error("Could not load configurations from file - trying environment configurations instead.")

		globalConfig, err = config.LoadEnv()
		if nil != err {
			log.WithError(err).Error("Could not load configurations from environment")
		}
	}
}

// initializes the basic configuration for the log wrapper
func initLog() {
	err := globalConfig.Log.Apply()
	if nil != err {
		log.WithError(err).Fatal("Could not apply logging configurations")
	}
}

func initStatsClient() {
	// FIXME: this causes application hang because we're in the locked log already
	//statsLog.SetHandler(func(msg string, fields map[string]interface{}, err error) {
	//	entry := log.WithFields(log.Fields(fields))
	//	if err == nil {
	//		entry.Warn(msg)
	//	} else {
	//		entry.WithError(err).Warn(msg)
	//	}
	//})

	sectionsTestsMap, err := bucket.ParseSectionsTestsMap(globalConfig.Stats.IDs)
	if err != nil {
		log.WithError(err).WithField("config", globalConfig.Stats.IDs).
			Error("Failed to parse stats second level IDs from env")
		sectionsTestsMap = map[bucket.PathSection]bucket.SectionTestDefinition{}
	}
	log.WithField("config", globalConfig.Stats.IDs).
		WithField("map", sectionsTestsMap.String()).
		Debug("Setting stats second level IDs")

	statsClient, err = stats.NewClient(globalConfig.Stats.DSN)
	if err != nil {
		log.WithError(err).Fatal("Error initializing stats client")
	}

	statsClient.SetHTTPMetricCallback(bucket.NewHasIDAtSecondLevelCallback(&bucket.SecondLevelIDConfig{
		HasIDAtSecondLevel:    sectionsTestsMap,
		AutoDiscoverThreshold: globalConfig.Stats.AutoDiscoverThreshold,
		AutoDiscoverWhiteList: globalConfig.Stats.AutoDiscoverWhiteList,
	}))

	host, err := os.Hostname()
	if nil != err {
		host = "-unknown-"
	}

	_, appFile := filepath.Split(os.Args[0])
	statsClient.TrackMetric("app", bucket.MetricOperation{"init", host, appFile})

	log.AddHook(hooks.NewLogrusHook(statsClient, globalConfig.Stats.ErrorsSection))
}

func initStatsExporter() {
	var err error
	logger := log.WithField("stats.exporter", globalConfig.Stats.Exporter)

	// Register stats exporter according to config
	switch globalConfig.Stats.Exporter {
	case "datadog":
	case "stackdriver":
		logger.Warn("Not implemented!")
		return
	case "prometheus":
		err = initPrometheusExporter()
		break
	default:
		logger.Info("Unsupported or invalid stats exporter was specified")
		return
	}

	if err != nil {
		logger.Warn("Failed initialising stats exporter")
		return
	}

	// Configure/Register stats views
	view.SetReportingPeriod(time.Second)

	vv := append(ochttp.DefaultServerViews, obs.AllViews...)

	if err := view.Register(vv...); err != nil {
		log.WithError(err).Warn("Failed to register server views")
	}
}

func initPrometheusExporter() (err error) {
	obs.PromExporter, err = prometheus.NewExporter(prometheus.Options{
		Namespace: globalConfig.Stats.Namespace,
	})
	if err != nil {
		log.WithError(err).Warn("Failed to create prometheus exporter")
	} else {
		view.RegisterExporter(obs.PromExporter)
	}
	return err
}

func initTracingExporter() {
	logger := log.WithField("tracing.exporter", globalConfig.Tracing.Exporter)

	switch globalConfig.Tracing.Exporter {
	case "azure_monitor":
	case "datadog":
	case "stackdriver":
	case "zipkin":
		logger.Warn("Not implemented!")
		return
	case "jaeger":
		initJaegerExporter()
		break
	default:
		logger.Warn("Unsupported or invalid tracing exporter was specified")
		return
	}
}

func initJaegerExporter() (err error) {
	jaegerExporter, err := jaeger.NewExporter(jaeger.Options{
		AgentEndpoint: globalConfig.Tracing.JaegerTracing.SamplingServerURL,
		ServiceName:   globalConfig.Tracing.ServiceName,
	})
	if err != nil {
		log.WithError(err).Warn("Failed to create jaeger exporter")
	} else {
		trace.RegisterExporter(jaegerExporter)
		trace.ApplyConfig(trace.Config{
			DefaultSampler: trace.ProbabilitySampler(globalConfig.Tracing.JaegerTracing.SamplingParam),
		})
	}
	return err
}
