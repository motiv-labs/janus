package main

import (
	"context"
	"net/url"
	"os"
	"strconv"
	"testing"

	"github.com/cucumber/godog"
	"github.com/stretchr/testify/assert"

	"github.com/hellofresh/janus/features/bootstrap"
	"github.com/hellofresh/janus/pkg/api"
	"github.com/hellofresh/janus/pkg/config"
)

const defaultUpstreamsPort = 9089

var (
	apiRepo          api.Repository
	cfg              *config.Specification
	cfgChan          chan api.ConfigurationMessage
	portSecondary    int
	apiPortSecondary int
	upstreamsPort    int
)

func InitializeTestSuite(ctx *godog.TestSuiteContext) {
	var err error

	ctx.BeforeSuite(func() {
		cfg, err = config.LoadEnv()
		if err != nil {
			panic(err)
		}

		dsnURL, err := url.Parse(cfg.Database.DSN)
		if nil != err {
			panic(err)
		}

		switch dsnURL.Scheme {
		case "mongodb":
			apiRepo, err = api.NewMongoAppRepository(cfg.Database.DSN, cfg.BackendFlushInterval)
			if err != nil {
				panic(err)
			}
		case "file":
			var apiPath = dsnURL.Path + "/apis"

			apiRepo, err = api.NewFileSystemRepository(apiPath)
			if err != nil {
				panic(err)
			}
		default:
			panic("invalid database")
		}

		portSecondary, err = strconv.Atoi(os.Getenv("PORT_SECONDARY"))
		if err != nil {
			panic(err)
		}

		apiPortSecondary, err = strconv.Atoi(os.Getenv("API_PORT_SECONDARY"))
		if err != nil {
			panic(err)
		}

		upstreamsPort = defaultUpstreamsPort
		if dynamicUpstreamsPort, exists := os.LookupEnv("DYNAMIC_UPSTREAMS_PORT"); exists {
			upstreamsPort, err = strconv.Atoi(dynamicUpstreamsPort)
			if err != nil {
				panic(err)
			}
		}

		cfgChan = make(chan api.ConfigurationMessage, 100)
		if listener, ok := apiRepo.(api.Listener); ok {
			listener.Listen(context.Background(), cfgChan)
		}
	})
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	bootstrap.RegisterRequestContext(ctx, cfg.Port, cfg.Web.Port, portSecondary, apiPortSecondary, defaultUpstreamsPort, upstreamsPort, cfg.Web.Credentials)
	bootstrap.RegisterAPIContext(ctx, apiRepo, cfgChan)
	bootstrap.RegisterMiscContext(ctx)
}

func Test_Fake(t *testing.T) {
	assert.True(t, true)
}
