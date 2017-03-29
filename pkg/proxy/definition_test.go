package proxy_test

import (
	"testing"

	"github.com/hellofresh/janus/pkg/proxy"
	"github.com/stretchr/testify/assert"
)

func TestSuccessfulValidation(t *testing.T) {
	definition := proxy.Definition{
		ListenPath: "/*",
	}

	assert.True(t, proxy.Validate(&definition))
}

func TestFailedValidation(t *testing.T) {
	definition := proxy.Definition{}

	assert.False(t, proxy.Validate(&definition))
}
