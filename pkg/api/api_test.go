package api_test

import (
	"testing"

	"github.com/hellofresh/janus/pkg/api"
	"github.com/stretchr/testify/assert"
)

const (
	testKey   = "key"
	testValue = "value"
)

func TestNewInstanceOfDefinition(t *testing.T) {
	instance := api.NewDefinition()

	assert.IsType(t, &api.Definition{}, instance)
	assert.True(t, instance.Active)
}
