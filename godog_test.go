package main

import (
	"context"
	"flag"
	"net/url"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/DATA-DOG/godog"
	"github.com/hellofresh/janus/features/bootstrap"
	"github.com/hellofresh/janus/pkg/api"
	"github.com/hellofresh/janus/pkg/config"
	"github.com/stretchr/testify/assert"
)

var (
	runGoDogTests bool
	stopOnFailure bool
)

func init() {
	flag.BoolVar(&runGoDogTests, "godog", false, "Set this flag is you want to run godog BDD tests")
	flag.BoolVar(&stopOnFailure, "stop-on-failure", false, "Stop processing on first failed scenario.. Flag is passed to godog")
	flag.Parse()
}

func FeatureContext(s *godog.Suite) {
	c, err := config.LoadEnv()
	if nil != err {
		panic(err)
	}

	var apiRepo api.Repository

	dsnURL, err := url.Parse(c.Database.DSN)
	switch dsnURL.Scheme {
	case "mongodb":
		apiRepo, err = api.NewMongoAppRepository(c.Database.DSN, c.BackendFlushInterval)
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

	portSecondary, err := strconv.Atoi(os.Getenv("PORT_SECONDARY"))
	if nil != err {
		panic(err)
	}

	apiPortSecondary, err := strconv.Atoi(os.Getenv("API_PORT_SECONDARY"))
	if nil != err {
		panic(err)
	}

	ch := make(chan api.ConfigurationMessage, 100)
	if listener, ok := apiRepo.(api.Listener); ok {
		listener.Listen(context.Background(), ch)
	}

	bootstrap.RegisterRequestContext(s, c.Port, c.Web.Port, portSecondary, apiPortSecondary, c.Web.Credentials)
	bootstrap.RegisterAPIContext(s, c.Web.ReadOnly, apiRepo, ch)
	bootstrap.RegisterMiscContext(s)
}

func Test_Fake(t *testing.T) {
	assert.True(t, true)
}

func TestMain(m *testing.M) {
	if !runGoDogTests {
		os.Exit(m.Run())
	}

	status := godog.RunWithOptions("Janus", func(s *godog.Suite) {
		FeatureContext(s)
	}, godog.Options{
		Format:        "pretty",
		Paths:         []string{"features"},
		Randomize:     time.Now().UTC().UnixNano(),
		StopOnFailure: stopOnFailure,
	})

	if st := m.Run(); st > status {
		status = st
	}

	os.Exit(status)
}
