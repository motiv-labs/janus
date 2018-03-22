package rate

import (
	"net/http"
	"testing"

	"github.com/hellofresh/janus/pkg/test"
	"github.com/hellofresh/stats-go"
	"github.com/stretchr/testify/assert"
	"github.com/ulule/limiter"
	smemory "github.com/ulule/limiter/drivers/store/memory"
)

func TestSuccessfulRateLimitLog(t *testing.T) {
	statsClient, _ := stats.NewClient("noop://")
	limiterStore := smemory.NewStore()
	rate, _ := limiter.NewRateFromFormatted("100-M")
	limiterInstance := limiter.New(limiterStore, rate)

	mw := NewRateLimitLogger(limiterInstance, statsClient)
	w, err := test.Record(
		"GET",
		"/",
		map[string]string{
			"Content-Type": "application/json",
		},
		mw(http.HandlerFunc(test.Ping)),
	)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
}
