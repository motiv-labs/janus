package compression

import (
	"testing"

	"github.com/hellofresh/janus/pkg/api"
	"github.com/hellofresh/janus/pkg/plugin"
	"github.com/hellofresh/janus/pkg/proxy"
	"github.com/stretchr/testify/assert"
)

func TestSetup(t *testing.T) {
	def := api.NewDefinition()
	route := proxy.NewRoute(def.Proxy)
	err := setupCompression(def, route, make(plugin.Config))
	assert.NoError(t, err)

	assert.Len(t, route.Inbound, 1)
}
