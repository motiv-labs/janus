package request_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/hellofresh/janus/pkg/request"
	"github.com/hellofresh/janus/pkg/test"
	"github.com/stretchr/testify/assert"
)

func TestBindSimpleJSON(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Body = ioutil.NopCloser(bytes.NewBuffer([]byte("{\"name\": \"Test Recipe\", \"tags\": [\"test\"]}")))

	recipe := test.Recipe{}
	err := request.BindJSON(req, &recipe)

	assert.Nil(t, err)
	assert.Equal(t, "Test Recipe", recipe.Name)
	assert.Equal(t, []test.Tag{"test"}, recipe.Tags)
}

func TestWrongBindJSON(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Body = ioutil.NopCloser(bytes.NewBuffer([]byte{}))

	recipe := test.Recipe{}
	err := request.BindJSON(req, &recipe)

	assert.NotNil(t, err)
}
