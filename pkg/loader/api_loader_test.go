// +build integration

package loader

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hellofresh/stats-go/client"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hellofresh/janus/pkg/api"
	"github.com/hellofresh/janus/pkg/errors"
	"github.com/hellofresh/janus/pkg/middleware"
	"github.com/hellofresh/janus/pkg/proxy"
	"github.com/hellofresh/janus/pkg/router"
	"github.com/hellofresh/janus/pkg/test"
)

const (
	defaultUpstreamsPort = "9089"
	defaultAPIsDir       = "../../assets/apis"
)

var tests = []struct {
	description     string
	method          string
	url             string
	headers         map[string]string
	expectedHeaders map[string]string
	expectedCode    int
}{
	{
		description: "Get example route",
		method:      "GET",
		url:         "/example",
		expectedHeaders: map[string]string{
			"Content-Type": "application/json; charset=utf-8",
		},
		expectedCode: http.StatusOK,
	}, {
		description: "Get invalid route",
		method:      "GET",
		url:         "/invalid-route",
		expectedHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		expectedCode: http.StatusNotFound,
	},
}

func TestSuccessfulLoader(t *testing.T) {
	log.SetOutput(ioutil.Discard)

	routerInstance := createRegisterAndRouter(t)
	ts := test.NewServer(routerInstance)
	defer ts.Close()

	for _, tc := range tests {
		res, err := ts.Do(tc.method, tc.url, tc.headers)
		require.NoError(t, err)
		t.Cleanup(func() {
			err := res.Body.Close()
			assert.NoError(t, err)
		})

		for headerName, headerValue := range tc.expectedHeaders {
			assert.Equal(t, headerValue, res.Header.Get(headerName))
		}

		assert.Equal(t, tc.expectedCode, res.StatusCode, tc.description)
	}
}

func createRegisterAndRouter(t *testing.T) router.Router {
	t.Helper()

	r := createRouter(t)
	r.Use(middleware.NewRecovery(errors.RecoveryHandler))

	register := proxy.NewRegister(proxy.WithRouter(r), proxy.WithStatsClient(client.NewNoop()))
	proxyRepo, err := api.NewFileSystemRepository(getAPIsDir(t))
	require.NoError(t, err)

	defs, err := proxyRepo.FindAll()
	require.NoError(t, err)

	loader := NewAPILoader(register)
	loader.RegisterAPIs(defs)

	return r
}

func createRouter(t *testing.T) router.Router {
	t.Helper()

	router.DefaultOptions.NotFoundHandler = errors.NotFound
	return router.NewChiRouterWithOptions(router.DefaultOptions)
}

func getAPIsDir(t *testing.T) string {
	t.Helper()

	upstreamsPort := os.Getenv("DYNAMIC_UPSTREAMS_PORT")
	if upstreamsPort == "" {
		// dynamic port is not set - use API defs as is
		return defaultAPIsDir
	}

	// dynamic port is set - we need to replace default port with the dynamic one
	dynamicAPIsDir, err := ioutil.TempDir("", "apis")
	require.NoError(t, err)
	t.Cleanup(func() {
		err := os.RemoveAll(dynamicAPIsDir)
		assert.NoError(t, err)
	})

	defaultUpstreamsHost := fmt.Sprintf("/localhost:%s/", defaultUpstreamsPort)
	dynamicUpstreamsHost := fmt.Sprintf("/localhost:%s/", upstreamsPort)
	err = filepath.Walk(defaultAPIsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}

		defaultContents, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		dynamicContents := strings.ReplaceAll(string(defaultContents), defaultUpstreamsHost, dynamicUpstreamsHost)
		return ioutil.WriteFile(filepath.Join(dynamicAPIsDir, info.Name()), []byte(dynamicContents), os.ModePerm)
	})
	require.NoError(t, err)

	return dynamicAPIsDir
}
