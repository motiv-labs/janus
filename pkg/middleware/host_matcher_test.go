package middleware_test

import (
	"testing"

	"net/http"

	"github.com/hellofresh/janus/pkg/middleware"
	"github.com/hellofresh/janus/pkg/test"
	"github.com/hellofresh/janus/pkg/web"
	"github.com/stretchr/testify/assert"
)

var (
	recovery = middleware.NewRecovery(web.RecoveryHandler).Handler
)

func TestMatchSimpleHeader(t *testing.T) {
	mw := middleware.NewHostMatcher([]string{"hellofresh.com"})
	w, err := test.Record(
		"GET",
		"/",
		map[string]string{
			"Content-Type": "application/json",
			"Host":         "hellofresh.com",
		},
		recovery(mw.Handler(http.HandlerFunc(test.Ping))),
	)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
}

func TestNotMatchSimpleHeader(t *testing.T) {
	mw := middleware.NewHostMatcher([]string{"hellofresh.com"})
	w, err := test.Record(
		"GET",
		"/",
		map[string]string{
			"Content-Type": "application/json",
			"Host":         "hellofresh.de",
		},
		recovery(mw.Handler(http.HandlerFunc(test.Ping))),
	)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
}

func TestMatchRegexHeader(t *testing.T) {
	mw := middleware.NewHostMatcher([]string{"hellofresh.*"})
	w, err := test.Record(
		"GET",
		"/",
		map[string]string{
			"Content-Type": "application/json",
			"Host":         "hellofresh.com",
		},
		recovery(mw.Handler(http.HandlerFunc(test.Ping))),
	)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
}

func TestNotMatchRegexHeader(t *testing.T) {
	mw := middleware.NewHostMatcher([]string{"hellofresh.*"})
	w, err := test.Record(
		"GET",
		"/",
		map[string]string{
			"Content-Type": "application/json",
			"Host":         "api.hellofresh.com",
		},
		recovery(mw.Handler(http.HandlerFunc(test.Ping))),
	)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
}
