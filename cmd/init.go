package cmd

import (
	"os"
	"path/filepath"

	"github.com/hellofresh/janus/pkg/config"
	"github.com/hellofresh/stats-go"
	"github.com/hellofresh/stats-go/bucket"
	"github.com/hellofresh/stats-go/client"
	"github.com/hellofresh/stats-go/hooks"
	log "github.com/sirupsen/logrus"
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

func initStatsd() {
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
		log.WithError(err).Fatal("Error initializing statsd client")
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
