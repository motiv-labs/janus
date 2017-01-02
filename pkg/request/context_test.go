package request_test

import (
	"testing"

	"github.com/hellofresh/janus/pkg/request"
	"github.com/stretchr/testify/assert"
)

func TestContextKey(t *testing.T) {
	key := request.ContextKey("test")
	assert.Equal(t, "janus.test", key.String())
}
