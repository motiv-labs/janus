package middleware

import (
	"net/http"
	"testing"

	"github.com/hellofresh/janus/pkg/test"
	"github.com/hellofresh/stats-go"
	"github.com/stretchr/testify/assert"
	"github.com/ulule/limiter"
)

func TestSuccessfulRateLimitLog(t *testing.T) {
	statsClient, _ := stats.NewClient("noop://", "")
	limiterStore := limiter.NewMemoryStore()
	rate, _ := limiter.NewRateFromFormatted("100-M")
	limiterInstance := limiter.NewLimiter(limiterStore, rate)

	mw := NewRateLimitLogger(limiterInstance, statsClient)
	w, err := test.Record(
		"GET",
		"/",
		map[string]string{
			"Content-Type": "application/json",
		},
		recovery(mw.Handler(http.HandlerFunc(test.Ping))),
	)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
}
