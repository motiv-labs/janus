package middleware

import (
	"testing"

	"net/http"

	"github.com/hellofresh/janus/pkg/test"
	"github.com/stretchr/testify/assert"
)

func TestSuccessfulRateLimitLog(t *testing.T) {
	mw := NewRateLimitLogger()
	w, err := test.Record(
		"GET",
		"/",
		map[string]string{
			"Content-Type": "application/json",
		},
		recovery(mw.Handler(http.HandlerFunc(test.Ping))),
	)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
}
