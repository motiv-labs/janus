package bodylmt

import (
	"testing"

	"github.com/hellofresh/janus/pkg/plugin"
	"github.com/hellofresh/janus/pkg/proxy"
	"github.com/stretchr/testify/assert"
)

func TestSetup(t *testing.T) {
	def := proxy.NewRouterDefinition(proxy.NewDefinition())
	err := setupBodyLimit(def, make(plugin.Config))
	assert.NoError(t, err)

	assert.Len(t, def.Middleware(), 1)
}
