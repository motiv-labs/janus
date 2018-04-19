package cb

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hellofresh/janus/pkg/api"
	"github.com/hellofresh/janus/pkg/test"
	"github.com/stretchr/testify/assert"
)

func TestMiddleware(t *testing.T) {
	t.Parallel()

	tests := []struct {
		scenario string
		function func(*testing.T, *http.Request, *httptest.ResponseRecorder)
	}{
		{
			scenario: "with wrong predicate given",
			function: testWrongPredicate,
		},
		{
			scenario: "when the upstream respond successfully",
			function: testSuccessfulUpstreamRetry,
		},
		{
			scenario: "when the upstream fails to respond",
			function: testFailedUpstreamRetry,
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()
			test.function(t, r, w)
		})
	}
}

func testWrongPredicate(t *testing.T, r *http.Request, w *httptest.ResponseRecorder) {
	cfg := Config{
		Predicate: "this is wrong",
	}
	d := api.Definition{
		Name: "example",
	}
	mw := NewCBMiddleware(cfg, &d)

	mw(http.HandlerFunc(test.Ping)).ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
}

func testSuccessfulUpstreamRetry(t *testing.T, r *http.Request, w *httptest.ResponseRecorder) {
	d := api.Definition{
		Name: "example",
	}
	mw := NewCBMiddleware(Config{}, &d)

	mw(http.HandlerFunc(test.Ping)).ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
}

func testFailedUpstreamRetry(t *testing.T, r *http.Request, w *httptest.ResponseRecorder) {
	d := api.Definition{
		Name: "example",
	}
	mw := NewCBMiddleware(Config{}, &d)

	mw(test.FailWith(http.StatusBadGateway)).ServeHTTP(w, r)

	assert.Equal(t, http.StatusBadGateway, w.Code)
}
