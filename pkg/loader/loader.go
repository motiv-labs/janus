package loader

import (
	"github.com/apex/log"
	httpErrors "github.com/hellofresh/janus/pkg/errors"
	"github.com/hellofresh/janus/pkg/plugin"
	"github.com/pkg/errors"
)

var loader *APILoader

func init() {
	plugin.RegisterEventHook(plugin.StartupEvent, onStartup)
	plugin.RegisterEventHook(plugin.ReloadEvent, onReload)
}

func onStartup(event interface{}) error {
	e, ok := event.(plugin.OnStartup)
	if !ok {
		return errors.New("Could not convert event to startup type")
	}

	loader = NewAPILoader(e.Register)
	loader.RegisterAPIs(e.Configuration)

	// some routers may panic when have empty routes list, so add one dummy 404 route to avoid this
	if e.Register.Router.RoutesCount() < 1 {
		e.Register.Router.Any("/", httpErrors.NotFound)
	}

	return nil
}

func onReload(event interface{}) error {
	e, ok := event.(plugin.OnReload)
	if !ok {
		return errors.New("Could not convert event to reload type")
	}

	if len(e.Configurations) == 0 {
		log.Debug("No configurations found to reload")
		return nil
	}

	loader.RegisterAPIs(e.Configurations)
	return nil
}
