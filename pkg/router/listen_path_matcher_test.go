package router_test

import (
	"testing"

	"github.com/hellofresh/janus/pkg/router"
	"github.com/stretchr/testify/assert"
)

var (
	matcher = router.NewListenPathMatcher()
)

func TestMatchCorrectRoute(t *testing.T) {
	assert.True(t, matcher.Match("/test/hello/*"))
	assert.True(t, matcher.Match("/test/hello/*/anything/after"))
}

func TestMatchIncorrectRoute(t *testing.T) {
	assert.False(t, matcher.Match("/test/hello"))
	assert.False(t, matcher.Match("/test/hello/anything/after"))
}

func TestExtractCorrectRoute(t *testing.T) {
	assert.Equal(t, "/test/hello", matcher.Extract("/test/hello/*"))
	assert.Equal(t, "/test/hello", matcher.Extract("/test/hello/*/anything/after"))
	assert.Equal(t, "/test/hello", matcher.Extract("/test/hello/*/anything/after/*"))
}

func TestExtractIncorrectRoute(t *testing.T) {
	assert.Equal(t, "/test/hello", matcher.Extract("/test/hello"))
	assert.Equal(t, "/test/hello/anything/after", matcher.Extract("/test/hello/anything/after"))
}
