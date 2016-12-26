package request_test

import (
	"testing"

	"github.com/hellofresh/janus/request"
	"github.com/stretchr/testify/assert"
)

// TestContextKey tests Rate methods.
func TestContextKey(t *testing.T) {
	key := request.ContextKey("test")
	assert.Equal(t, "janus.test", key.String())
}