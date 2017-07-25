package bodylmt

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hellofresh/janus/pkg/test"
	"github.com/stretchr/testify/assert"
)

func TestBodyLmtValidSize(t *testing.T) {
	mw := NewBodyLimitMiddleware("2M")

	content := []byte("Hello, World!")
	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(content))
	w := httptest.NewRecorder()

	mw(http.HandlerFunc(test.Ping)).ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestBodyLmtInvalidSize(t *testing.T) {
	mw := NewBodyLimitMiddleware("2B")

	content := []byte("Hello, World!")
	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(content))
	w := httptest.NewRecorder()

	mw(http.HandlerFunc(test.Ping)).ServeHTTP(w, r)

	assert.Equal(t, http.StatusRequestEntityTooLarge, w.Code)
}
