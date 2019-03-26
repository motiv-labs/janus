package basic

import (
	"errors"

	"github.com/hellofresh/janus/pkg/plugin"
	"github.com/hellofresh/janus/pkg/proxy"
	"github.com/hellofresh/janus/pkg/router"
	log "github.com/sirupsen/logrus"
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

func setupBasicAuth(def *proxy.RouterDefinition, rawConfig plugin.Config) error {
	if repo == nil {
		return errors.New("the repository was not set by onStartup event")
	}

	def.AddMiddleware(NewBasicAuth(repo))
	return nil
}

func onAdminAPIStartup(event interface{}) error {
	e, ok := event.(plugin.OnAdminAPIStartup)
	if !ok {
		return errors.New("could not convert event to admin startup type")
	}

	adminRouter = e.Router
	return nil
}

func onStartup(event interface{}) error {
	var err error

	e, ok := event.(plugin.OnStartup)
	if !ok {
		return errors.New("could not convert event to startup type")
	}

	if e.MongoSession == nil {
		log.Debug("Mongo session is not set, using memory repository for basic auth plugin")

		repo = NewInMemoryRepository()
	} else {
		log.Debug("Mongo session is set, using mongo repository for basic auth plugin")

		repo, err = NewMongoRepository(e.MongoSession)
		if err != nil {
			return err
		}
	}

	if adminRouter == nil {
		return ErrInvalidAdminRouter
	}

	handlers := NewHandler(repo)
	group := adminRouter.Group("/credentials/basic_auth")
	{
		group.GET("/", handlers.Index())
		group.POST("/", handlers.Create())
		group.GET("/{username}", handlers.Show())
		group.PUT("/{username}", handlers.Update())
		group.DELETE("/{username}", handlers.Delete())
	}

	return nil
}
