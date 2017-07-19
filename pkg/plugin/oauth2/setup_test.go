package oauth2

import (
	"encoding/json"
	"testing"

	"github.com/hellofresh/janus/pkg/oauth"
	"github.com/hellofresh/janus/pkg/plugin"
	"github.com/hellofresh/janus/pkg/proxy"
	"github.com/hellofresh/janus/pkg/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOAuth2Config(t *testing.T) {
	var config Config
	rawConfig := map[string]interface{}{
		"server_name": "test",
	}

	err := plugin.Decode(rawConfig, &config)
	require.NoError(t, err)
	assert.Equal(t, "test", config.ServerName)
}

func TestSetupWithValidOAuthServer(t *testing.T) {
	rawConfig := map[string]interface{}{
		"server_name": "test",
	}
	tokenStrategyRaw, _ := json.Marshal(map[string]string{"secret": "1234"})

	repo := oauth.NewInMemoryRepository()
	repo.Add(&oauth.OAuth{
		Name: "test",
		TokenStrategy: oauth.TokenStrategy{
			Name:     "jwt",
			Settings: oauth.TokenStrategySettings(tokenStrategyRaw),
		},
	})

	route := proxy.NewRoute(&proxy.Definition{})
	err := setupOAuth2(route, plugin.Params{
		Config:    rawConfig,
		Storage:   store.NewInMemoryStore(),
		OAuthRepo: repo,
	})

	require.NoError(t, err)
	assert.Len(t, route.Inbound, 2)
}

func TestSetupWithInvalidOAuthServer(t *testing.T) {
	rawConfig := map[string]interface{}{
		"server_name": "test",
	}
	tokenStrategyRaw, _ := json.Marshal(map[string]string{"secret": "1234"})

	repo := oauth.NewInMemoryRepository()
	repo.Add(&oauth.OAuth{
		Name: "test1",
		TokenStrategy: oauth.TokenStrategy{
			Name:     "jwt",
			Settings: oauth.TokenStrategySettings(tokenStrategyRaw),
		},
	})
	route := proxy.NewRoute(&proxy.Definition{})
	err := setupOAuth2(route, plugin.Params{
		Config:    rawConfig,
		Storage:   store.NewInMemoryStore(),
		OAuthRepo: repo,
	})

	assert.Error(t, err)
}

func TestSetupWithWrongStrategy(t *testing.T) {
	rawConfig := map[string]interface{}{
		"server_name": "test",
	}
	tokenStrategyRaw, _ := json.Marshal(map[string]string{"secret": "1234"})

	repo := oauth.NewInMemoryRepository()
	repo.Add(&oauth.OAuth{
		Name: "test1",
		TokenStrategy: oauth.TokenStrategy{
			Name:     "wrong",
			Settings: oauth.TokenStrategySettings(tokenStrategyRaw),
		},
	})
	route := proxy.NewRoute(&proxy.Definition{})
	err := setupOAuth2(route, plugin.Params{
		Config:    rawConfig,
		Storage:   store.NewInMemoryStore(),
		OAuthRepo: repo,
	})

	assert.Error(t, err)
}
