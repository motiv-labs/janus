package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hellofresh/janus/pkg/test"
	"github.com/stretchr/testify/assert"
)

func TestAddHeader(t *testing.T) {
	config := RequestTransformerConfig{
		Add: RequestTransformerOptions{
			Headers: map[string]string{
				"Test": "Test",
			},
		},
	}
	mw := NewRequestTransformer(config)
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	mw.Handler(http.HandlerFunc(test.Ping)).ServeHTTP(w, req)
	assert.Equal(t, "Test", req.Header.Get("Test"))
}

func TestAddHeaderThatAlreadyExists(t *testing.T) {
	config := RequestTransformerConfig{
		Add: RequestTransformerOptions{
			Headers: map[string]string{
				"Test": "New value",
			},
		},
	}
	mw := NewRequestTransformer(config)
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("Test", "Original value")
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	mw.Handler(http.HandlerFunc(test.Ping)).ServeHTTP(w, req)
	assert.Equal(t, "Original value", req.Header.Get("Test"))
}

func TestAppendHeader(t *testing.T) {
	config := RequestTransformerConfig{
		Append: RequestTransformerOptions{
			Headers: map[string]string{
				"Test": "New value",
			},
		},
	}
	mw := NewRequestTransformer(config)
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("Test2", "Original value")
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	mw.Handler(http.HandlerFunc(test.Ping)).ServeHTTP(w, req)
	assert.Equal(t, "New value", req.Header.Get("Test"))
	assert.Equal(t, "Original value", req.Header.Get("Test2"))
}

func TestReplaceHeader(t *testing.T) {
	config := RequestTransformerConfig{
		Replace: RequestTransformerOptions{
			Headers: map[string]string{
				"Test": "New value",
			},
		},
	}
	mw := NewRequestTransformer(config)
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("Test", "Original value")
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	mw.Handler(http.HandlerFunc(test.Ping)).ServeHTTP(w, req)
	assert.Equal(t, "New value", req.Header.Get("Test"))
}

func TestReplaceHeaderThatDoesntExist(t *testing.T) {
	config := RequestTransformerConfig{
		Replace: RequestTransformerOptions{
			Headers: map[string]string{
				"Test": "New value",
			},
		},
	}
	mw := NewRequestTransformer(config)
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	mw.Handler(http.HandlerFunc(test.Ping)).ServeHTTP(w, req)
	assert.Equal(t, "", req.Header.Get("Test"))
}

func TestRemoveHeaderThatDoesntExist(t *testing.T) {
	config := RequestTransformerConfig{
		Remove: RequestTransformerOptions{
			Headers: map[string]string{
				"Test": "New value",
			},
		},
	}
	mw := NewRequestTransformer(config)
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	mw.Handler(http.HandlerFunc(test.Ping)).ServeHTTP(w, req)
	assert.Equal(t, "", req.Header.Get("Test"))
}

func TestRemoveHeader(t *testing.T) {
	config := RequestTransformerConfig{
		Remove: RequestTransformerOptions{
			Headers: map[string]string{
				"Test": "New value",
			},
		},
	}
	mw := NewRequestTransformer(config)
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("Test", "Original value")
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	mw.Handler(http.HandlerFunc(test.Ping)).ServeHTTP(w, req)
	assert.Equal(t, "", req.Header.Get("Test"))
}
