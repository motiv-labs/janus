package compression

import (
	"testing"

	"github.com/hellofresh/janus/pkg/plugin"
	"github.com/hellofresh/janus/pkg/proxy"
	"github.com/stretchr/testify/assert"
)

func TestSetup(t *testing.T) {
	route := proxy.NewRoute(&proxy.Definition{})
	err := setupCompression(route, plugin.Params{})
	assert.NoError(t, err)

	assert.Len(t, route.Inbound, 1)
}
