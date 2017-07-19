package oauth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTokenStrategySettings_GetJWTSigningMethods(t *testing.T) {
	settingsLegacy := TokenStrategySettings(`{"secret": "foo-bar"}`)
	methodsLegacy, err := settingsLegacy.GetJWTSigningMethods()
	require.NoError(t, err)

	require.Equal(t, 1, len(methodsLegacy))
	assert.Equal(t, "HS256", methodsLegacy[0].Alg)
	assert.Equal(t, "foo-bar", methodsLegacy[0].Key)

	settingsList := TokenStrategySettings(`[{"alg":"HS256","key":"foo"},{"alg":"RS256","key":"bar"}]`)
	methodsList, err := settingsList.GetJWTSigningMethods()
	require.NoError(t, err)

	require.Equal(t, 2, len(methodsList))
	assert.Equal(t, "HS256", methodsList[0].Alg)
	assert.Equal(t, "foo", methodsList[0].Key)
	assert.Equal(t, "RS256", methodsList[1].Alg)
	assert.Equal(t, "bar", methodsList[1].Key)

	settingsLegacyError := TokenStrategySettings(`{"foo": "bar"}`)
	_, err = settingsLegacyError.GetJWTSigningMethods()
	assert.Equal(t, ErrJWTSecretMissing, err)

	settingsInvalid := TokenStrategySettings(`"`)
	_, err = settingsInvalid.GetJWTSigningMethods()
	assert.Error(t, err)
}
