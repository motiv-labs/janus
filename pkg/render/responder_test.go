package render_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hellofresh/janus/pkg/render"
	"github.com/hellofresh/janus/pkg/test"
	"github.com/stretchr/testify/assert"
)

func TestRespondAsJson(t *testing.T) {
	w := httptest.NewRecorder()

	recipe := test.Recipe{Name: "Test"}
	render.JSON(w, http.StatusOK, recipe)

	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRespondExpectedBody(t *testing.T) {
	w := httptest.NewRecorder()

	recipe := test.Recipe{Name: "Test"}
	render.JSON(w, http.StatusOK, recipe)

	expectedWriter := httptest.NewRecorder()
	json.NewEncoder(expectedWriter).Encode(recipe)

	assert.Equal(t, expectedWriter.Body, w.Body)
}

func TestWrongJson(t *testing.T) {
	w := httptest.NewRecorder()

	render.JSON(w, http.StatusOK, make(chan int))

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
