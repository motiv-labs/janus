package checker

import (
	"testing"

	"net/http"

	"github.com/hellofresh/janus/pkg/api"
	"github.com/hellofresh/janus/pkg/router"
	"github.com/hellofresh/janus/pkg/test"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	r := router.NewChiRouter()
	repo := api.NewInMemoryRepository()

	err := Register(r, repo)
	assert.NoError(t, err)

	ts := test.NewServer(r)
	defer ts.Close()

	res, err := ts.Do(http.MethodGet, "/status", make(map[string]string))
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "application/json", res.Header.Get("Content-Type"))
}
