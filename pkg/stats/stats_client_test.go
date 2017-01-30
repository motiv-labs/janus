package stats

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatsMetricName(t *testing.T) {
	client := &StatsClient{nil}

	dataProvider := []struct {
		Method string
		Path   string
		Metric string
	}{
		{"GET", "/", "get.-.-"},
		{"TRACE", "/api", "trace.api.-"},
		{"TRACE", "/api/", "trace.api.-"},
		{"POST", "/api/recipes", "post.api.recipes"},
		{"POST", "/api/recipes/", "post.api.recipes"},
		{"DELETE", "/api/recipes/123", "delete.api.recipes"},
		{"DELETE", "/api/recipes.foo-bar/123", "delete.api.recipes_foo-bar"},
		{"DELETE", "/api/recipes.foo_bar/123", "delete.api.recipes_foo__bar"},
	}

	for _, data := range dataProvider {
		assert.Equal(t, data.Metric, client.getStatsdMetricName(data.Method, &url.URL{Path: data.Path}))
	}
}
