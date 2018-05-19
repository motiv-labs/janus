package responsetransformer

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hellofresh/janus/pkg/test"
	"github.com/stretchr/testify/assert"
)

func TestAddHeader(t *testing.T) {
	config := Config{
		Add: Options{
			Headers: map[string]string{
				"Test": "Test",
			},
		},
	}

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	NewResponseTransformer(config)(http.HandlerFunc(test.Ping)).ServeHTTP(w, req)

	assert.Equal(t, "Test", w.Header().Get("Test"))
}

func TestReplaceHeader(t *testing.T) {
	config := Config{
		Replace: Options{
			Headers: map[string]string{
				"Content-Type": "test",
			},
		},
	}

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	NewResponseTransformer(config)(http.HandlerFunc(test.Ping)).ServeHTTP(w, req)

	assert.Equal(t, "test", w.Header().Get("Content-Type"))
}

func TestRemoveHeader(t *testing.T) {
	config := Config{
		Remove: Options{
			Headers: map[string]string{
				"Content-Type": "",
			},
		},
	}

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	NewResponseTransformer(config)(http.HandlerFunc(test.Ping)).ServeHTTP(w, req)

	assert.Equal(t, "", w.Header().Get("Content-Type"))
}

func TestAppendHeader(t *testing.T) {
	config := Config{
		Append: Options{
			Headers: map[string]string{
				"Test": "test",
			},
		},
	}

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	NewResponseTransformer(config)(http.HandlerFunc(test.Ping)).ServeHTTP(w, req)

	assert.Equal(t, "test", w.Header().Get("Test"))
}
