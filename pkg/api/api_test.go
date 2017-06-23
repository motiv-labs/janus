package api_test

import (
	"testing"

	"github.com/hellofresh/janus/pkg/api"
	"github.com/stretchr/testify/assert"
)

func TestNewInstanceOfDefinition(t *testing.T) {
	instance := api.NewDefinition()

	assert.IsType(t, &api.Definition{}, instance)
	assert.True(t, instance.Active)
}

func TestSuccessfulValidation(t *testing.T) {
	instance := api.NewDefinition()
	instance.Name = "Test"
	isValid, err := instance.Validate()

	assert.NoError(t, err)
	assert.True(t, isValid)
}

func TestFailedValidation(t *testing.T) {
	instance := api.NewDefinition()
	isValid, err := instance.Validate()

	assert.Error(t, err)
	assert.False(t, isValid)
}
