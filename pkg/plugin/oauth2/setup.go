package oauth2

import (
	"github.com/hellofresh/janus/pkg/jwt"
	"github.com/hellofresh/janus/pkg/oauth"
	"github.com/hellofresh/janus/pkg/plugin"
	"github.com/hellofresh/janus/pkg/proxy"
	log "github.com/sirupsen/logrus"
)

func init() {
	plugin.RegisterPlugin("oauth2", plugin.Plugin{
		Action: setupOAuth2,
	})
}

// Config represents the oauth configuration
type Config struct {
	ServerName string `json:"server_name"`
}

func setupOAuth2(route *proxy.Route, p plugin.Params) error {
	var config Config
	err := plugin.Decode(p.Config, &config)
	if err != nil {
		return err
	}

	oauthServer, err := p.OAuthRepo.FindByName(config.ServerName)
	if nil != err {
		return err
	}

	manager, err := getManager(oauthServer, config.ServerName)
	if nil != err {
		log.WithError(err).Error("OAuth Configuration for this API is incorrect, skipping...")
		return err
	}

	secret, err := oauthServer.TokenStrategy.Settings.GetJWTSecret()
	if err != nil {
		return err
	}

	route.AddInbound(NewKeyExistsMiddleware(manager))
	route.AddInbound(NewRevokeRulesMiddleware(jwt.NewParser(jwt.NewConfig(secret)), oauthServer.AccessRules))

	return nil
}

func getManager(oauthServer *oauth.OAuth, oAuthServerName string) (oauth.Manager, error) {
	managerType, err := oauth.ParseType(oauthServer.TokenStrategy.Name)
	if nil != err {
		return nil, err
	}

	return oauth.NewManagerFactory(oauthServer).Build(managerType)
}
