package basic

import (
	"testing"

	"gopkg.in/mgo.v2"

	"github.com/hellofresh/janus/pkg/plugin"
	"github.com/hellofresh/janus/pkg/proxy"
	"github.com/hellofresh/janus/pkg/router"
	"github.com/stretchr/testify/require"
)

func TestSetup(t *testing.T) {
	route := proxy.NewRoute(&proxy.Definition{})

	event1 := plugin.OnAdminAPIStartup{Router: router.NewChiRouter()}
	err := onAdminAPIStartup(event1)
	require.NoError(t, err)

	event2 := plugin.OnStartup{Register: proxy.NewRegister(router.NewChiRouter(), proxy.Params{}), MongoSession: &mgo.Session{}}
	err = onStartup(event2)
	require.NoError(t, err)

	err = setupBasicAuth(route, make(plugin.Config))
	require.NoError(t, err)
}

func TestOnStartupMissingMongoSession(t *testing.T) {
	event := plugin.OnStartup{Register: proxy.NewRegister(router.NewChiRouter(), proxy.Params{})}
	err := onStartup(event)
	require.Error(t, err)
	require.IsType(t, ErrInvalidMongoDBSession, err)
}

func TestOnStartupMissingAdminRouter(t *testing.T) {
	event := plugin.OnStartup{}
	err := onStartup(event)
	require.Error(t, err)
	require.IsType(t, ErrInvalidAdminRouter, err)
}

func TestOnStartupWrongEvent(t *testing.T) {
	wrongEvent := plugin.OnAdminAPIStartup{}
	err := onStartup(wrongEvent)
	require.Error(t, err)
}

func TestOnAdminAPIStartupWrongEvent(t *testing.T) {
	wrongEvent := plugin.OnStartup{}
	err := onAdminAPIStartup(wrongEvent)
	require.Error(t, err)
}
