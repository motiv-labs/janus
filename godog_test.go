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

func FeatureContext(s *godog.Suite) {
	c, err := config.LoadEnv()
	if nil != err {
		panic(err)
	}

	var apiRepo api.Repository

	dsnURL, err := url.Parse(c.Database.DSN)
	if nil != err {
		panic(err)
	}

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
	bootstrap.RegisterAPIContext(s, apiRepo, ch)
	bootstrap.RegisterMiscContext(s)
}

func Test_Fake(t *testing.T) {
	assert.True(t, true)
}
