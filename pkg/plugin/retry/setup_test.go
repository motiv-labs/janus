package retry

import (
	"testing"

	"github.com/hellofresh/janus/pkg/api"
	"github.com/hellofresh/janus/pkg/plugin"
	"github.com/stretchr/testify/assert"
)

func TestSetup(t *testing.T) {
	def := api.NewDefinition()
	err := setupRetry(def, make(plugin.Config))
	assert.NoError(t, err)

	assert.Len(t, def.Proxy.Middleware(), 1)
}
