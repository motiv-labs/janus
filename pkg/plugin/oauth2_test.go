package plugin

import (
	"testing"

	"github.com/hellofresh/janus/pkg/api"
	"github.com/hellofresh/janus/pkg/oauth"
	"github.com/hellofresh/janus/pkg/store"
	"github.com/stretchr/testify/assert"
)

func TestOAuth2Config(t *testing.T) {
	var config oauth2Config
	rawConfig := map[string]interface{}{
		"server_name": "test",
	}

	err := decode(rawConfig, &config)
	assert.NoError(t, err)
	assert.Equal(t, "test", config.ServerName)
}

func TestOAuth2PluginGetName(t *testing.T) {
	plugin := NewOAuth2(oauth.NewInMemoryRepository(), store.NewInMemoryStore())
	assert.Equal(t, "oauth2", plugin.GetName())
}

func TestOAtuh2PluginWithValidOAuthServer(t *testing.T) {
	rawConfig := map[string]interface{}{
		"server_name": "test",
	}

	spec := &api.Spec{
		Definition: &api.Definition{
			Name: "API Name",
		},
	}

	repo := oauth.NewInMemoryRepository()
	repo.Add(&oauth.OAuth{
		Name: "test",
		TokenStrategy: oauth.TokenStrategy{
			Name:     "jwt",
			Settings: oauth.TokenStrategySettings{"secret": "1234"},
		},
	})
	plugin := NewOAuth2(repo, store.NewInMemoryStore())
	middleware, err := plugin.GetMiddlewares(rawConfig, spec)

	assert.NoError(t, err)
	assert.Len(t, middleware, 1)
}

func TestOAtuh2PluginWithInalidOAuthServer(t *testing.T) {
	rawConfig := map[string]interface{}{
		"server_name": "test",
	}

	spec := &api.Spec{
		Definition: &api.Definition{
			Name: "API Name",
		},
	}

	repo := oauth.NewInMemoryRepository()
	repo.Add(&oauth.OAuth{
		Name: "test1",
		TokenStrategy: oauth.TokenStrategy{
			Name:     "jwt",
			Settings: oauth.TokenStrategySettings{"secret": "1234"},
		},
	})
	plugin := NewOAuth2(repo, store.NewInMemoryStore())
	_, err := plugin.GetMiddlewares(rawConfig, spec)

	assert.Error(t, err)
}

func TestOAtuh2PluginWithWrongStrategy(t *testing.T) {
	rawConfig := map[string]interface{}{
		"server_name": "test",
	}

	spec := &api.Spec{
		Definition: &api.Definition{
			Name: "API Name",
		},
	}

	repo := oauth.NewInMemoryRepository()
	repo.Add(&oauth.OAuth{
		Name: "test1",
		TokenStrategy: oauth.TokenStrategy{
			Name:     "wrong",
			Settings: oauth.TokenStrategySettings{"secret": "1234"},
		},
	})
	plugin := NewOAuth2(repo, store.NewInMemoryStore())
	_, err := plugin.GetMiddlewares(rawConfig, spec)

	assert.Error(t, err)
}
