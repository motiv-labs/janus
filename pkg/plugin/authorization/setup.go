package authorization

import (
	"errors"

	"github.com/hellofresh/janus/pkg/config"
	"github.com/hellofresh/janus/pkg/models"
	"github.com/hellofresh/janus/pkg/plugin"
	"github.com/hellofresh/janus/pkg/proxy"
)

var (
	tokenManager *models.TokenManager
	roleManager  *models.RoleManager
)

func init() {
	err := plugin.RegisterEventHook(plugin.StartupEvent, onStartup)
	if err != nil {
		panic(err)
	}

	err = plugin.RegisterPlugin("authorization", plugin.Plugin{
		Action:   setupAuthorization,
		Validate: nil,
	})
	if err != nil {
		panic(err)
	}
}

func setupAuthorization(def *proxy.RouterDefinition, _ plugin.Config) error {
	def.AddMiddleware(NewTokenCheckerMiddleware())
	def.AddMiddleware(NewRoleCheckerMiddleware())

	return nil
}

func onStartup(event interface{}) error {
	var err error
	var conf config.Config

	_, ok := event.(plugin.OnStartup)
	if !ok {
		return errors.New("could not convert event to startup type")
	}

	err = config.UnmarshalYAML("./config/config.yaml", &conf)
	if err != nil {
		return err
	}
	conf.KafkaConfig.Normalize()

	if err != nil {
		return err
	}

	tokenManager = &models.TokenManager{Tokens: map[string]*models.JWTToken{}}
	roleManager = &models.RoleManager{Roles: map[string]*models.Role{}}

	err = RefreshTokens(&conf, tokenManager)
	if err != nil && !errors.Is(err, ErrTimeout) {
		return err
	}
	err = RefreshRoles(&conf, roleManager)
	if err != nil && !errors.Is(err, ErrTimeout) {
		return err
	}

	StartFactConsumers(&conf, tokenManager, roleManager)

	return nil
}
