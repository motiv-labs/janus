package rate

import (
	"context"
	"net"
	"net/http"

	"github.com/felixge/httpsnoop"
	"github.com/hellofresh/stats-go/bucket"
	"github.com/hellofresh/stats-go/client"
	log "github.com/sirupsen/logrus"
	"github.com/ulule/limiter"
)

const (
	limiterSection = "limiter"
	limiterMetric  = "state"
)

// NewRateLimitLogger logs the IP of blocked users with rate limit
func NewRateLimitLogger(lmt *limiter.Limiter, statsClient client.Client, trustForwardHeaders bool) func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Debug("Starting RateLimitLogger.WriterWrapper middleware")

			m := httpsnoop.CaptureMetrics(handler, w, r)

			limiterIP := limiter.GetIP(r, trustForwardHeaders)
			if m.Code == http.StatusTooManyRequests {
				log.WithFields(log.Fields{
					"ip_address":  limiterIP.String(),
					"request_uri": r.RequestURI,
				}).Warning("Rate Limit exceeded for this IP")
			}

			trackLimitState(lmt, statsClient, limiterIP, r)
		})
	}
}

func trackLimitState(lmt *limiter.Limiter, statsClient client.Client, limiterIP net.IP, r *http.Request) {
	context, err := lmt.Peek(context.Background(), limiterIP.String())
	if err != nil {
		log.WithError(err).WithFields(log.Fields{
			"ip_address":  limiterIP.String(),
			"request_uri": r.RequestURI,
		}).Error("Failed to get limiter context from request")
	} else {
		requestsPerformed := context.Limit - context.Remaining
		limitState := requestsPerformed * 100 / context.Limit

		operation := bucket.BuildHTTPRequestMetricOperation(r, statsClient.GetHTTPMetricCallback())
		// replace request method with fixed section name
		operation[0] = limiterMetric

		statsClient.TrackState(limiterSection, operation, int(limitState))
	}
}
