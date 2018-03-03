package middleware

import (
	"net/http"
	"testing"

	"github.com/hellofresh/janus/pkg/test"
	"github.com/hellofresh/stats-go"
	"github.com/hellofresh/stats-go/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSuccessfulStats(t *testing.T) {
	statsClient, err := stats.NewClient("memory://")
	require.NoError(t, err)

	mw := NewStats(statsClient)
	w, err := test.Record(
		http.MethodGet,
		"/",
		map[string]string{
			"Content-Type": "application/json",
		},
		mw.Handler(http.HandlerFunc(test.Ping)),
	)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	require.IsType(t, &client.Memory{}, statsClient)
	memoryClient := statsClient.(*client.Memory)

	require.Equal(t, 1, len(memoryClient.TimerMetrics), 1)
	assert.Equal(t, "request.get.-.-", memoryClient.TimerMetrics[0].Bucket)

	require.Equal(t, 4, len(memoryClient.CountMetrics))
	assert.Equal(t, 1, memoryClient.CountMetrics["request.get.-.-"])
	assert.Equal(t, 1, memoryClient.CountMetrics["request-ok.get.-.-"])
	assert.Equal(t, 1, memoryClient.CountMetrics["total.request"])
	assert.Equal(t, 1, memoryClient.CountMetrics["total.request-ok"])
}

func TestUnknownPath(t *testing.T) {
	statsClient, err := stats.NewClient("memory://")
	require.NoError(t, err)

	mw := NewStats(statsClient)
	w, err := test.Record(
		http.MethodGet,
		"/",
		map[string]string{
			"Content-Type": "application/json",
		},
		mw.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(http.StatusText(http.StatusNotFound)))
		})),
	)
	require.NoError(t, err)

	assert.Equal(t, http.StatusNotFound, w.Code)

	require.IsType(t, &client.Memory{}, statsClient)
	memoryClient := statsClient.(*client.Memory)

	require.Equal(t, 1, len(memoryClient.TimerMetrics), 1)
	assert.Equal(t, "request.get.-not-found-.-", memoryClient.TimerMetrics[0].Bucket)

	require.Equal(t, 4, len(memoryClient.CountMetrics))
	require.Equal(t, 4, len(memoryClient.CountMetrics))
	assert.Equal(t, 1, memoryClient.CountMetrics["request.get.-not-found-.-"])
	assert.Equal(t, 1, memoryClient.CountMetrics["request-fail.get.-not-found-.-"])
	assert.Equal(t, 1, memoryClient.CountMetrics["total.request"])
	assert.Equal(t, 1, memoryClient.CountMetrics["total.request-fail"])
}
