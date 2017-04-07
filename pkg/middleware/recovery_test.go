package middleware

import (
	"testing"

	"net/http"

	"github.com/hellofresh/janus/pkg/errors"
	"github.com/hellofresh/janus/pkg/test"
	"github.com/stretchr/testify/assert"
)

func TestSuccessfulRecovery(t *testing.T) {
	mw := NewRecovery(test.RecoveryHandler)
	w, err := test.Record(
		"GET",
		"/",
		map[string]string{
			"Content-Type": "application/json",
		},
		mw.Handler(http.HandlerFunc(doPanic)),
	)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
}

func doPanic(w http.ResponseWriter, r *http.Request) {
	panic(errors.ErrInvalidID)
}
