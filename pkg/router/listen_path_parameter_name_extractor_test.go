package router_test

import (
	"testing"

	"github.com/hellofresh/janus/pkg/router"
	"github.com/stretchr/testify/assert"
)

func TestExtractCorrectParameters(t *testing.T) {
	extractor := router.NewListenPathParamNameExtractor()

	assert.Equal(t, []string{}, extractor.Extract("/recipes/"))
	assert.Equal(t, []string{}, extractor.Extract("/recipes?take=100"))
	assert.Equal(t, []string{"id"}, extractor.Extract("/recipes/{id}/favorites"))
	assert.Equal(t, []string{"id", "slug"}, extractor.Extract("/recipes/{id}/favorites/{slug}"))
	assert.Equal(t, []string{"id", "slug"}, extractor.Extract("/recipes/{id}/favorites/{slug}?q=123"))
}
