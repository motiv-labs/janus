package oauth2

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/asaskevich/govalidator"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"

	"github.com/hellofresh/janus/pkg/config"
	"github.com/hellofresh/janus/pkg/jwt"
	"github.com/hellofresh/janus/pkg/plugin"
	"github.com/hellofresh/janus/pkg/proxy"
	"github.com/hellofresh/janus/pkg/router"
)

const (
	mongodb   = "mongodb"
	file      = "file"
	cassandra = "cassandra"

	mongoIdxTimeout = 10 * time.Second
)

var (
	repo        Repository
	loader      *OAuthLoader
	adminRouter router.Router
)

func init() {
	plugin.RegisterEventHook(plugin.StartupEvent, onStartup)
	plugin.RegisterEventHook(plugin.ReloadEvent, onReload)
	plugin.RegisterEventHook(plugin.AdminAPIStartupEvent, onAdminAPIStartup)
	plugin.RegisterPlugin("oauth2", plugin.Plugin{
		Action:   setupOAuth2,
		Validate: validateConfig,
	})
}

// Config represents the oauth configuration
type Config struct {
	ServerName string `json:"server_name"`
}

func onAdminAPIStartup(event interface{}) error {
	e, ok := event.(plugin.OnAdminAPIStartup)
	if !ok {
		return errors.New("could not convert event to admin startup type")
	}

	adminRouter = e.Router
	return nil
}

func onReload(event interface{}) error {
	_, ok := event.(plugin.OnReload)
	if !ok {
		return errors.New("could not convert event to reload type")
	}

	loader.LoadDefinitions(repo)

	return nil
}

func onStartup(event interface{}) error {
	e, ok := event.(plugin.OnStartup)
	if !ok {
		return errors.New("could not convert event to startup type")
	}

	cfg := e.Config.Database
	dsnURL, err := url.Parse(cfg.DSN)
	if err != nil {
		return err
	}

	switch dsnURL.Scheme {
	case mongodb:
		repo, err = NewMongoRepository(e.MongoDB)
		if err != nil {
			return fmt.Errorf("could not create a mongodb repository for oauth servers: %w", err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), mongoIdxTimeout)
		defer cancel()

		if _, err := e.MongoDB.Collection(collectionName).Indexes().CreateOne(
			ctx,
			mongo.IndexModel{
				Keys: bsonx.Doc{
					{Key: "name", Value: bsonx.Int32(1)},
				},
				Options: options.Index().SetUnique(true).SetBackground(true).SetSparse(true),
			},
		); err != nil {
			return fmt.Errorf("failed to create indexes for oauth servers repository: %w", err)
		}
	case cassandra:
		repo, err = NewCassandraRepository(e.Cassandra)
		if err != nil {
			log.Errorf("error creating new cassandra repo")
			return err
		}
	case file:
		authPath := fmt.Sprintf("%s/auth", dsnURL.Path)
		log.WithField("path", authPath).Debug("Trying to load Auth configuration files")

		repo, err = NewFileSystemRepository(authPath)
		if err != nil {
			return fmt.Errorf("could not create a file based repository for the oauth servers: %w", err)
		}

	default:
		return errors.New("the selected scheme is not supported to load OAuth servers")
	}

	loadOAuthEndpoints(adminRouter, repo, e.Config.Web.Credentials)
	loader = NewOAuthLoader(e.Register)
	loader.LoadDefinitions(repo)

	return nil
}

func setupOAuth2(def *proxy.RouterDefinition, rawConfig plugin.Config) error {
	var pluginConf Config
	err := plugin.Decode(rawConfig, &pluginConf)
	if err != nil {
		return err
	}

	oauthServer, err := repo.FindByName(pluginConf.ServerName)
	if nil != err {
		return err
	}

	manager, err := getManager(oauthServer, pluginConf.ServerName)
	if nil != err {
		log.WithError(err).Error("OAuth Configuration for this API is incorrect, skipping...")
		return err
	}
	signingMethods, err := oauthServer.TokenStrategy.GetJWTSigningMethods()
	if err != nil {
		return err
	}

	def.AddMiddleware(NewKeyExistsMiddleware(manager))
	def.AddMiddleware(NewRevokeRulesMiddleware(jwt.NewParser(jwt.NewParserConfig(oauthServer.TokenStrategy.Leeway, signingMethods...)), oauthServer.AccessRules))

	return nil
}

func validateConfig(rawConfig plugin.Config) (bool, error) {
	var cfg Config
	err := plugin.Decode(rawConfig, &cfg)
	if err != nil {
		return false, err
	}

	return govalidator.ValidateStruct(cfg)
}

func getManager(oauthServer *OAuth, oAuthServerName string) (Manager, error) {
	managerType, err := ParseType(oauthServer.TokenStrategy.Name)
	if nil != err {
		return nil, err
	}

	return NewManagerFactory(oauthServer).Build(managerType)
}

// loadOAuthEndpoints register api endpoints
func loadOAuthEndpoints(router router.Router, repo Repository, cred config.Credentials) {
	log.Debug("Loading OAuth Endpoints")

	guard := jwt.NewGuard(cred)
	oAuthHandler := NewController(repo)
	oauthGroup := router.Group("/oauth/servers")
	oauthGroup.Use(jwt.NewMiddleware(guard).Handler)
	{
		oauthGroup.GET("/", oAuthHandler.Get())
		oauthGroup.GET("/{name}", oAuthHandler.GetBy())
		oauthGroup.POST("/", oAuthHandler.Post())
		oauthGroup.PUT("/{name}", oAuthHandler.PutBy())
		oauthGroup.DELETE("/{name}", oAuthHandler.DeleteBy())
	}
}
