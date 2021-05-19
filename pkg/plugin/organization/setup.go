package organization

import (
	"errors"
	"github.com/hellofresh/janus/pkg/plugin"
	"github.com/hellofresh/janus/pkg/plugin/basic"
	"github.com/hellofresh/janus/pkg/proxy"
	"github.com/hellofresh/janus/pkg/router"
	log "github.com/sirupsen/logrus"
)

var (
	repo        Repository
	basicRepo basic.Repository
	adminRouter router.Router
)

// Organization represents the configuration to save the user and organization pair
type Organization struct {
	Username string `json:"username"`
	Organization  string `json:"organization"`
	Password string `json:"password"`
}

func init() {
	plugin.RegisterEventHook(plugin.StartupEvent, onStartup)
	plugin.RegisterEventHook(plugin.AdminAPIStartupEvent, onAdminAPIStartup)

	plugin.RegisterPlugin("organization_auth", plugin.Plugin{
		Action: setupOrganization,
	})
}

func setupOrganization(def *proxy.RouterDefinition, rawConfig plugin.Config) error {
	if repo == nil {
		return errors.New("the repository was not set by onStartup event")
	}

	var organization Organization
	err := plugin.Decode(rawConfig, &organization)
	if err != nil {
		return err
	}

	def.AddMiddleware(NewOrganization(organization, repo))
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

	if e.MongoDB != nil {
		log.Debug("Mongo DB is set, using mongo repository for organization plugin")
		log.Debug("unimplemented")
		//repo, err = NewMongoRepository(e.MongoDB)
		if err != nil {
			return err
		}
	} else if e.Cassandra != nil {
		log.Debugf("Cassandra is set, using cassandra repository for organization plugin")

		repo, err = NewCassandraRepository(e.Cassandra)
		if err != nil {
			log.Errorf("error getting cassandra repo: %v", err)
			return err
		}
	} else {
		log.Debug("No DB set, using memory repository for organization plugin")
		log.Debug("unimplemented")
		//repo = NewInMemoryRepository()
	}

	if adminRouter == nil {
		return ErrInvalidAdminRouter
	}

	handlers := NewHandler(repo)
	group := adminRouter.Group("/credentials/organization_auth")
	{
		group.GET("/", handlers.Index())
		group.POST("/", handlers.Create())
		group.GET("/{username}", handlers.Show())
		group.PUT("/{username}", handlers.Update())
		group.DELETE("/{username}", handlers.Delete())
	}

	return nil
}
