package responsetransformer

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddHeader(t *testing.T) {
	config := Config{
		Add: Options{
			Headers: map[string]string{
				"Test": "Test",
			},
		},
	}

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"fake twitter json string"}`)
	}))

	resp, err := http.Get(s.URL)
	require.NoError(t, err)
	NewResponseTransformer(config)(resp.Request, resp)

	assert.Equal(t, "Test", resp.Header.Get("Test"))
}

func TestReplaceHeader(t *testing.T) {
	config := Config{
		Replace: Options{
			Headers: map[string]string{
				"Content-Type": "test",
			},
		},
	}

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"fake twitter json string"}`)
	}))

	resp, err := http.Get(s.URL)
	require.NoError(t, err)
	NewResponseTransformer(config)(resp.Request, resp)

	assert.Equal(t, "test", resp.Header.Get("Content-Type"))
}

func TestRemoveHeader(t *testing.T) {
	config := Config{
		Remove: Options{
			Headers: map[string]string{
				"Content-Type": "",
			},
		},
	}

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"fake twitter json string"}`)
	}))

	resp, err := http.Get(s.URL)
	require.NoError(t, err)
	NewResponseTransformer(config)(resp.Request, resp)

	assert.Equal(t, "", resp.Header.Get("Content-Type"))
}

func TestAppendHeader(t *testing.T) {
	config := Config{
		Append: Options{
			Headers: map[string]string{
				"Test": "test",
			},
		},
	}

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"fake twitter json string"}`)
	}))

	resp, err := http.Get(s.URL)
	require.NoError(t, err)
	NewResponseTransformer(config)(resp.Request, resp)

	assert.Equal(t, "test", resp.Header.Get("Test"))
}
