package request_test

import (
	"bytes"
	"io/ioutil"
	"testing"

	"net/http"

	"github.com/hellofresh/janus/mock"
	"github.com/hellofresh/janus/request"
	"github.com/stretchr/testify/assert"
)

func TestBindSimpleJson(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Body = ioutil.NopCloser(bytes.NewBuffer([]byte("{\"name\": \"Test Recipe\", \"tags\": [\"test\"]}")))

	recipe := mock.Recipe{}
	err := request.BindJSON(req, &recipe)

	assert.Nil(t, err)
	assert.Equal(t, "Test Recipe", recipe.Name)
	assert.Equal(t, []mock.Tag{"test"}, recipe.Tags)
}
