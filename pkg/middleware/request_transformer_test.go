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

func TestAddQueryString(t *testing.T) {
	config := RequestTransformerConfig{
		Add: RequestTransformerOptions{
			QueryString: map[string]string{
				"test": "value",
			},
		},
	}
	mw := NewRequestTransformer(config)
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	mw.Handler(http.HandlerFunc(test.Ping)).ServeHTTP(w, req)
	assert.Equal(t, "value", req.URL.Query().Get("test"))
}

func TestAddQueryStringThatAlreadyExists(t *testing.T) {
	config := RequestTransformerConfig{
		Add: RequestTransformerOptions{
			QueryString: map[string]string{
				"test": "new-value",
			},
		},
	}
	mw := NewRequestTransformer(config)
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	q := req.URL.Query()
	q.Add("test", "original-value")
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	mw.Handler(http.HandlerFunc(test.Ping)).ServeHTTP(w, req)
	assert.Equal(t, "original-value", q.Get("test"))
}

func TestAppendQueryString(t *testing.T) {
	config := RequestTransformerConfig{
		Append: RequestTransformerOptions{
			QueryString: map[string]string{
				"test": "value",
			},
		},
	}
	mw := NewRequestTransformer(config)
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	mw.Handler(http.HandlerFunc(test.Ping)).ServeHTTP(w, req)

	q := req.URL.Query()
	q.Add("test2", "value")

	assert.Equal(t, "value", q.Get("test"))
	assert.Equal(t, "value", q.Get("test2"))
}

func TestReplaceQueryString(t *testing.T) {
	config := RequestTransformerConfig{
		Replace: RequestTransformerOptions{
			QueryString: map[string]string{
				"test": "new-value",
			},
		},
	}
	mw := NewRequestTransformer(config)
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.NoError(t, err)

	q := req.URL.Query()
	q.Add("test", "original-value")
	req.URL.RawQuery = q.Encode()

	w := httptest.NewRecorder()
	mw.Handler(http.HandlerFunc(test.Ping)).ServeHTTP(w, req)
	assert.Equal(t, "new-value", req.URL.Query().Get("test"))
}

func TestReplaceQueryStringThatDoesntExists(t *testing.T) {
	config := RequestTransformerConfig{
		Replace: RequestTransformerOptions{
			QueryString: map[string]string{
				"test": "new-value",
			},
		},
	}
	mw := NewRequestTransformer(config)
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	mw.Handler(http.HandlerFunc(test.Ping)).ServeHTTP(w, req)
	assert.Equal(t, "", req.URL.Query().Get("test"))
}

func TestRemoveQueryString(t *testing.T) {
	config := RequestTransformerConfig{
		Remove: RequestTransformerOptions{
			QueryString: map[string]string{
				"test": "",
			},
		},
	}
	mw := NewRequestTransformer(config)
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.NoError(t, err)

	q := req.URL.Query()
	q.Add("test", "original-value")
	req.URL.RawQuery = q.Encode()

	w := httptest.NewRecorder()
	mw.Handler(http.HandlerFunc(test.Ping)).ServeHTTP(w, req)
	assert.Equal(t, "", req.URL.Query().Get("test"))
}

func TestRemoveQueryStringthatDoesntExists(t *testing.T) {
	config := RequestTransformerConfig{
		Remove: RequestTransformerOptions{
			QueryString: map[string]string{
				"test": "",
			},
		},
	}
	mw := NewRequestTransformer(config)
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	mw.Handler(http.HandlerFunc(test.Ping)).ServeHTTP(w, req)
	assert.Equal(t, "", req.URL.Query().Get("test"))
}
