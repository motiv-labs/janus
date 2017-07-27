package errors

import (
	errorsBase "errors"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorWithCustomError(t *testing.T) {
	w := httptest.NewRecorder()
	Handler(w, New(http.StatusBadRequest, http.StatusText(http.StatusBadRequest)))

	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestErrorWithDefaultError(t *testing.T) {
	w := httptest.NewRecorder()
	Handler(w, errorsBase.New(http.StatusText(http.StatusBadRequest)))

	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
