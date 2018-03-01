package loader

import (
	"github.com/hellofresh/janus/pkg/api"
	httpErrors "github.com/hellofresh/janus/pkg/errors"
	"github.com/hellofresh/janus/pkg/plugin"
	"github.com/hellofresh/janus/pkg/proxy"
	"github.com/pkg/errors"
)

func init() {
	plugin.RegisterEventHook(plugin.StartupEvent, onStartup)
	plugin.RegisterEventHook(plugin.ReloadEvent, onReload)
}

func onStartup(event interface{}) error {
	e, ok := event.(plugin.OnStartup)
	if !ok {
		return errors.New("Could not convert event to startup type")
	}

	Load(e.Register, e.Repository)
	return nil
}

func onReload(event interface{}) error {
	e, ok := event.(plugin.OnReload)
	if !ok {
		return errors.New("Could not convert event to reload type")
	}

	Load(e.Register, e.Repository)
	return nil
}

// Load loads all the basic components and definitions into a router
func Load(register *proxy.Register, repo api.Repository) {
	apiLoader := NewAPILoader(register)
	apiLoader.LoadDefinitions(repo)

	// some routers may panic when have empty routes list, so add one dummy 404 route to avoid this
	if register.Router.RoutesCount() < 1 {
		register.Router.Any("/", httpErrors.NotFound)
	}
}
