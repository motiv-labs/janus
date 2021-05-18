package company

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

// Company represents the configuration to save the user and company pair
type Company struct {
	Username string `json:"username"`
	Company  string `json:"company"`
}

func init() {
	log.Debugf("in init for company setup")
	plugin.RegisterEventHook(plugin.StartupEvent, onStartup)
	plugin.RegisterEventHook(plugin.AdminAPIStartupEvent, onAdminAPIStartup)

	plugin.RegisterPlugin("company_plugin", plugin.Plugin{
		Action: setupCompany,
	})
}

func setupCompany(def *proxy.RouterDefinition, rawConfig plugin.Config) error {
	if repo == nil {
		return errors.New("the repository was not set by onStartup event")
	}

	var company Company
	err := plugin.Decode(rawConfig, &company)
	if err != nil {
		return err
	}

	def.AddMiddleware(NewCompany(company, repo))
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
		log.Debug("Mongo DB is set, using mongo repository for company plugin")
		log.Debug("unimplemented")
		//repo, err = NewMongoRepository(e.MongoDB)
		if err != nil {
			return err
		}
	} else if e.Cassandra != nil {
		log.Debugf("Cassandra is set, using cassandra repository for company plugin")

		repo, err = NewCassandraRepository(e.Cassandra)
		if err != nil {
			log.Errorf("error getting cassandra repo: %v", err)
			return err
		}
	} else {
		log.Debug("No DB set, using memory repository for company plugin")
		log.Debug("unimplemented")
		//repo = NewInMemoryRepository()
	}

	if adminRouter == nil {
		return ErrInvalidAdminRouter
	}

	handlers := NewHandler(repo)
	group := adminRouter.Group("/credentials/company")
	{
		group.GET("/", handlers.Index())
		group.POST("/", handlers.Create())
		group.GET("/{username}", handlers.Show())
		group.PUT("/{username}", handlers.Update())
		group.DELETE("/{username}", handlers.Delete())
	}

	return nil
}
