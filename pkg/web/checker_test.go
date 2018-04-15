package web

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

	r.GET("/status", NewOverviewHandler(repo))

	ts := test.NewServer(r)
	defer ts.Close()

	res, _ := ts.Do(http.MethodGet, "/status", make(map[string]string))
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "application/json", res.Header.Get("Content-Type"))
}
