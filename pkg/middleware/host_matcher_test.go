package middleware

import (
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/hellofresh/janus/pkg/test"
	"github.com/stretchr/testify/assert"
)

var (
	recovery = NewRecovery(test.RecoveryHandler).Handler
)

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}

func TestMatchSimpleHeader(t *testing.T) {
	mw := NewHostMatcher([]string{"hellofresh.com"})
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
	mw := NewHostMatcher([]string{"hellofresh.com"})
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
	mw := NewHostMatcher([]string{"hellofresh.*"})
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
	mw := NewHostMatcher([]string{"hellofresh.*"})
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
