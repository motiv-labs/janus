package basic

import (
	"testing"

	"github.com/globalsign/mgo"

	"github.com/hellofresh/janus/pkg/plugin"
	"github.com/hellofresh/janus/pkg/proxy"
	"github.com/hellofresh/janus/pkg/router"
	"github.com/stretchr/testify/require"
)

func TestSetup(t *testing.T) {
	def := proxy.NewRouterDefinition(proxy.NewDefinition())

	event1 := plugin.OnAdminAPIStartup{Router: router.NewChiRouter()}
	err := onAdminAPIStartup(event1)
	require.NoError(t, err)

	event2 := plugin.OnStartup{Register: proxy.NewRegister(proxy.WithRouter(router.NewChiRouter())), MongoSession: &mgo.Session{}}
	err = onStartup(event2)
	require.NoError(t, err)

	err = setupBasicAuth(def, make(plugin.Config))
	require.NoError(t, err)
}

func TestOnStartupMissingAdminRouter(t *testing.T) {
	// reset admin router to avoid dependency from another test
	adminRouter = nil

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
