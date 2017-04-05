package response_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hellofresh/janus/pkg/response"
	"github.com/hellofresh/janus/pkg/test"
	"github.com/stretchr/testify/assert"
)

func TestRespondAsJson(t *testing.T) {
	w := httptest.NewRecorder()

	recipe := test.Recipe{Name: "Test"}
	response.JSON(w, http.StatusOK, recipe)

	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Equal(t, w.Code, http.StatusOK)
}

func TestRespondExpectedBody(t *testing.T) {
	w := httptest.NewRecorder()

	recipe := test.Recipe{Name: "Test"}
	response.JSON(w, http.StatusOK, recipe)

	expectedWriter := httptest.NewRecorder()
	json.NewEncoder(expectedWriter).Encode(recipe)

	assert.Equal(t, expectedWriter.Body, w.Body)
}
