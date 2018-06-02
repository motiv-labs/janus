package oauth2

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/mgo.v2/bson"
)

func TestAccessRulesWithWrongPredicate(t *testing.T) {
	rules := AccessRule{
		Action:    "wrong",
		Predicate: "wrong.predicate",
	}
	_, err := rules.IsAllowed(make(map[string]interface{}))
	require.Error(t, err)
}

func TestAccessRulesWithEmptyPredicate(t *testing.T) {
	rules := AccessRule{
		Action:    "wrong",
		Predicate: "",
	}
	_, err := rules.IsAllowed(map[string]interface{}{"test": true})
	require.Error(t, err)
}

func TestAccessRulesWithPredicateThatDoesntMatch(t *testing.T) {
	rules := AccessRule{
		Action:    "wrong",
		Predicate: "test = false",
	}
	_, err := rules.IsAllowed(map[string]interface{}{"test": true})
	require.Error(t, err)
}

func TestTokenStrategyWithInvalidSettings(t *testing.T) {
	settingsLegacy := TokenStrategy{Settings: make(chan int)}
	_, err := settingsLegacy.GetJWTSigningMethods()
	require.Error(t, err)
}

func TestTokenStrategyWithEmptySecret(t *testing.T) {
	settingsLegacy := TokenStrategy{Settings: bson.M{"secret": ""}}
	_, err := settingsLegacy.GetJWTSigningMethods()
	require.Error(t, err)
}

func TestTokenStrategy_GetJWTSigningMethods_mongo(t *testing.T) {
	settingsLegacy := TokenStrategy{Settings: bson.M{"secret": "foo-bar"}}
	methodsLegacy, err := settingsLegacy.GetJWTSigningMethods()
	require.NoError(t, err)

	require.Equal(t, 1, len(methodsLegacy))
	assert.Equal(t, "HS256", methodsLegacy[0].Alg)
	assert.Equal(t, "foo-bar", methodsLegacy[0].Key)

	settingsList := TokenStrategy{Settings: []interface{}{bson.M{"alg": "HS256", "key": "foo"}, bson.M{"alg": "RS256", "key": "bar"}}}
	methodsList, err := settingsList.GetJWTSigningMethods()
	require.NoError(t, err)

	require.Equal(t, 2, len(methodsList))
	assert.Equal(t, "HS256", methodsList[0].Alg)
	assert.Equal(t, "foo", methodsList[0].Key)
	assert.Equal(t, "RS256", methodsList[1].Alg)
	assert.Equal(t, "bar", methodsList[1].Key)

	settingsLegacyError := TokenStrategy{Settings: bson.M{"foo": "bar"}}
	_, err = settingsLegacyError.GetJWTSigningMethods()
	assert.Equal(t, ErrJWTSecretMissing, err)

	settingsInvalid := TokenStrategy{Settings: bson.M{}}
	_, err = settingsInvalid.GetJWTSigningMethods()
	assert.Error(t, err)
}

func TestTokenStrategy_GetJWTSigningMethods_file(t *testing.T) {
	repo := &FileSystemRepository{}

	oauthServerLegacy := repo.parseOAuthServer([]byte(fileLegacy))
	methodsLegacy, err := oauthServerLegacy.TokenStrategy.GetJWTSigningMethods()
	require.NoError(t, err)

	require.Equal(t, 1, len(methodsLegacy))
	assert.Equal(t, "HS256", methodsLegacy[0].Alg)
	assert.Equal(t, "foo-bar", methodsLegacy[0].Key)

	oauthServerList := repo.parseOAuthServer([]byte(fileList))

	methodsList, err := oauthServerList.TokenStrategy.GetJWTSigningMethods()
	require.NoError(t, err)

	require.Equal(t, 2, len(methodsList))
	assert.Equal(t, "HS256", methodsList[0].Alg)
	assert.Equal(t, "foo", methodsList[0].Key)
	assert.Equal(t, "RS256", methodsList[1].Alg)
	assert.Equal(t, "bar", methodsList[1].Key)

	oauthServerLegacyError := repo.parseOAuthServer([]byte(fileLegacyError))
	_, err = oauthServerLegacyError.TokenStrategy.GetJWTSigningMethods()
	assert.Equal(t, ErrJWTSecretMissing, err)

	oauthServerInvalid := repo.parseOAuthServer([]byte(fileBad))
	_, err = oauthServerInvalid.TokenStrategy.GetJWTSigningMethods()
	assert.Error(t, err)
}

const (
	fileLegacy = `{
    "name" : "legacy",
    "oauth_endpoints" : {
        "token" : {
            "preserve_host" : false,
            "listen_path" : "/auth/token",
            "upstreams": {"balancing": "roundrobin", "targets": [{"target": "http://localhost:8080/token"}]},
            "strip_path" : true,
            "append_path" : false,
            "methods" : [
                "GET",
                "POST"
            ]
        }
    },
    "token_strategy" : {
        "name" : "jwt",
        "settings" : {"secret": "foo-bar"}
    }
}`
	fileList = `{
    "name" : "list",
    "oauth_endpoints" : {
        "token" : {
            "preserve_host" : false,
            "listen_path" : "/auth/token",
            "upstreams": {"balancing": "roundrobin", "targets": [{"target": "http://localhost:8080/token"}]},
            "strip_path" : true,
            "append_path" : false,
            "methods" : [
                "GET",
                "POST"
            ]
        }
    },
    "token_strategy" : {
        "name" : "jwt",
        "settings" : [
            {"alg": "HS256", "key": "foo"},
            {"alg": "RS256", "key": "bar"}
        ]
    }
}`
	fileLegacyError = `{
    "name" : "legacy",
    "oauth_endpoints" : {
        "token" : {
            "preserve_host" : false,
            "listen_path" : "/auth/token",
            "upstreams": {"balancing": "roundrobin", "targets": [{"target": "http://localhost:8080/token"}]},
            "strip_path" : true,
            "append_path" : false,
            "methods" : [
                "GET",
                "POST"
            ]
        }
    },
    "token_strategy" : {
        "name" : "jwt",
        "settings" : {"foo": "bar"}
    }
}`
	fileBad = `{
    "name" : "legacy",
    "oauth_endpoints" : {
        "token" : {
            "preserve_host" : false,
            "listen_path" : "/auth/token",
            "upstreams": {"balancing": "roundrobin", "targets": [{"target": "http://localhost:8080/token"}]},
            "strip_path" : true,
            "append_path" : false,
            "methods" : [
                "GET",
                "POST"
            ]
        }
    },
    "token_strategy" : {
        "name" : "jwt",
        "settings" : {}
    }
}`
)
