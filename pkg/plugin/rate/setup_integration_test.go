// +build integration

package rate

import (
	"testing"

	"github.com/hellofresh/janus/pkg/api"
	"github.com/stretchr/testify/assert"
)

func TestRateLimitPluginRedisPolicy(t *testing.T) {
	rawConfig := map[string]interface{}{
		"limit":  "10-S",
		"policy": "redis",
		"redis": map[string]interface{}{
			"dsn":    "localhost",
			"prefix": "test",
		},
	}

	def := api.NewDefinition()
	err := setupRateLimit(def, rawConfig)

	assert.Error(t, err)
}
