package basic

import (
	"errors"

	"github.com/hellofresh/janus/pkg/plugin"
	"github.com/hellofresh/janus/pkg/proxy"
	"github.com/hellofresh/janus/pkg/router"
)

var (
	repo        Repository
	adminRouter router.Router
)

func init() {
	plugin.RegisterEventHook(plugin.StartupEvent, onStartup)
	plugin.RegisterEventHook(plugin.AdminAPIStartupEvent, onAdminAPIStartup)

	plugin.RegisterPlugin("basic_auth", plugin.Plugin{
		Action: setupBasicAuth,
	})
}

func setupBasicAuth(route *proxy.Route, rawConfig plugin.Config) error {
	if repo == nil {
		return errors.New("The repository was not set by onStartup event")
	}

	route.AddInbound(NewBasicAuth(repo))
	return nil
}

func onAdminAPIStartup(event interface{}) error {
	e, ok := event.(plugin.OnAdminAPIStartup)
	if !ok {
		return errors.New("Could not convert event to admin startup type")
	}

	adminRouter = e.Router
	return nil
}

func onStartup(event interface{}) error {
	var err error

	e, ok := event.(plugin.OnStartup)
	if !ok {
		return errors.New("Could not convert event to startup type")
	}

	if e.MongoSession == nil {
		return ErrInvalidMongoDBSession
	}

	if adminRouter == nil {
		return ErrInvalidAdminRouter
	}

	repo, err = NewMongoRepository(e.MongoSession)
	if err != nil {
		return err
	}

	handlers := NewController(repo)
	group := adminRouter.Group("/credentials/basic_auth")
	{
		group.GET("/", handlers.Get())
		group.POST("/", handlers.Post())
		group.GET("/{username}", handlers.GetBy())
		group.PUT("/{username}", handlers.PutBy())
		group.DELETE("/{username}", handlers.DeleteBy())
	}

	return nil
}
