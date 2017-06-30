package oauth2

import (
	"github.com/hellofresh/janus/pkg/oauth"
	"github.com/hellofresh/janus/pkg/plugin"
	"github.com/hellofresh/janus/pkg/proxy"
	"github.com/hellofresh/janus/pkg/store"
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

	manager, err := getManager(p.OAuthRepo, p.Storage, config.ServerName)
	if nil != err {
		log.WithError(err).Error("OAuth Configuration for this API is incorrect, skipping...")
		return err
	}
	route.AddInbound(NewKeyExistsMiddleware(manager))

	return nil
}

func getManager(authRepo oauth.Repository, storage store.Store, oAuthServerName string) (oauth.Manager, error) {
	oauthServer, err := authRepo.FindByName(oAuthServerName)
	if nil != err {
		return nil, err
	}

	managerType, err := oauth.ParseType(oauthServer.TokenStrategy.Name)
	if nil != err {
		return nil, err
	}

	return oauth.NewManagerFactory(storage, oauthServer.TokenStrategy.Settings).Build(managerType)
}
