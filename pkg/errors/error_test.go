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
	r, _ := http.NewRequest(http.MethodGet, "/hello/test", nil)
	Handler(w, r, New(http.StatusBadRequest, http.StatusText(http.StatusBadRequest)))

	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestErrorWithDefaultError(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/hello/test", nil)
	Handler(w, r, errorsBase.New(http.StatusText(http.StatusBadRequest)))

	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestErrorNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	NotFound(w, r)

	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRecovery(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	RecoveryHandler(w, r, ErrInvalidID)

	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
