package plugin

import (
	"github.com/hellofresh/janus/pkg/api"
	"github.com/hellofresh/janus/pkg/middleware"
	"github.com/hellofresh/janus/pkg/oauth"
	"github.com/hellofresh/janus/pkg/router"
	"github.com/hellofresh/janus/pkg/store"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
)

type oauth2Config struct {
	ServerName string `json:"server_name"`
}

// OAuth2 checks the integrity of the provided OAuth headers
type OAuth2 struct {
	authRepo oauth.Repository
	storage  store.Store
}

// NewOAuth2 creates a new instance of KeyExistsMiddleware
func NewOAuth2(authRepo oauth.Repository, storage store.Store) *OAuth2 {
	return &OAuth2{authRepo, storage}
}

// GetName retrieves the plugin's name
func (h *OAuth2) GetName() string {
	return "oauth2"
}

// GetMiddlewares retrieves the plugin's middlewares
func (h *OAuth2) GetMiddlewares(rawConfig map[string]interface{}, referenceSpec *api.Spec) ([]router.Constructor, error) {
	var oauth2Config oauth2Config
	err := mapstructure.Decode(rawConfig, &oauth2Config)
	if err != nil {
		return nil, err
	}

	manager, err := h.getManager(oauth2Config.ServerName)
	if nil != err {
		log.WithError(err).Error("OAuth Configuration for this API is incorrect, skipping...")
		return nil, err
	}

	mw := middleware.NewKeyExistsMiddleware(manager)
	return []router.Constructor{
		mw.Handler,
	}, nil
}

func (h *OAuth2) getManager(oAuthServerName string) (oauth.Manager, error) {
	oauthServer, err := h.authRepo.FindByName(oAuthServerName)
	if nil != err {
		return nil, err
	}

	managerType, err := oauth.ParseType(oauthServer.TokenStrategy.Name)
	if nil != err {
		return nil, err
	}

	return oauth.NewManagerFactory(h.storage, oauthServer.TokenStrategy.Settings).Build(managerType)
}
