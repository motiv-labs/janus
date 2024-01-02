package authorization

import (
	"github.com/hellofresh/janus/pkg/config"
	"github.com/hellofresh/janus/pkg/plugin"
	"github.com/hellofresh/janus/pkg/proxy"
)

func init() {
	err := plugin.RegisterPlugin("authorization", plugin.Plugin{
		Action:   setupAuthorization,
		Validate: nil,
	})
	if err != nil {
		panic(err)
	}
}

func setupAuthorization(def *proxy.RouterDefinition, _ plugin.Config) error {
	var conf config.Config

	err := config.UnmarshalYAML("./config/config.yaml", &conf)
	if err != nil {
		return err
	}

	tokenManager := NewTokenManager(&conf)
	roleManager := NewRoleManager(&conf)

	err = tokenManager.FetchTokens()
	if err != nil {
		return err
	}
	err = roleManager.FetchRoles()
	if err != nil {
		return err
	}

	// Catch middleware needs to be first, if it successfully catches – it will interrupt http request.
	def.AddMiddleware(NewTokenCatcherMiddleware(tokenManager))
	def.AddMiddleware(NewRoleCatcherMiddleware(roleManager))

	def.AddMiddleware(NewTokenCheckerMiddleware(tokenManager))
	def.AddMiddleware(NewRoleCheckerMiddleware(roleManager))

	return nil
}
