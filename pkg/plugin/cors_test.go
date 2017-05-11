package plugin

import (
	"testing"

	"github.com/hellofresh/janus/pkg/api"
	"github.com/stretchr/testify/assert"
)

func TestCORS_GetMiddlewares_convertToSlice(t *testing.T) {
	config := api.Config{"domains": []interface{}{"*"}}
	slice := convertToSlice(config["domains"])
	assert.Equal(t, []string{"*"}, slice)

	config = api.Config{"domains": []interface{}{"api.example.com", "gui.example.com"}}
	slice = convertToSlice(config["domains"])
	assert.Equal(t, []string{"api.example.com", "gui.example.com"}, slice)

	unknown := convertToSlice(config["methods"])
	assert.Equal(t, []string{}, unknown)
}
