package middleware

import (
	"net"
	"net/http"
	"sync"

	"github.com/hellofresh/janus/pkg/response"
	"github.com/hellofresh/stats-go"
	"github.com/hellofresh/stats-go/bucket"
	log "github.com/sirupsen/logrus"
	"github.com/ulule/limiter"
)

const (
	limiterSection = "limiter"
	limiterMetric  = "state"
)

// RateLimitLogger represents the middleware struct
type RateLimitLogger struct {
	limiter     *limiter.Limiter
	statsClient stats.Client
}

// NewRateLimitLogger logs the IP of blocked users with rate limit
func NewRateLimitLogger(limiter *limiter.Limiter, statsClient stats.Client) *RateLimitLogger {
	return &RateLimitLogger{limiter, statsClient}
}

// Handler is the middleware function
func (m *RateLimitLogger) Handler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			lock          sync.Mutex
			headerWritten bool
		)

		log.Debug("Starting RateLimitLogger.WriterWrapper middleware")

		hooks := response.Hooks{
			WriteHeader: func(next response.WriteHeaderFunc) response.WriteHeaderFunc {
				return func(code int) {
					next(code)
					lock.Lock()
					defer lock.Unlock()
					if !headerWritten {
						limiterIP := limiter.GetIP(r)
						if code == http.StatusTooManyRequests {
							log.WithFields(log.Fields{
								"ip_address":  limiterIP.String(),
								"request_uri": r.RequestURI,
							}).Warning("Rate Limit exceded for this IP")
						}

						m.trackLimitState(limiterIP, r)

						headerWritten = true
					}
				}
			},
		}

		handler.ServeHTTP(response.Wrap(w, hooks), r)
	})
}

func (m *RateLimitLogger) trackLimitState(limiterIP net.IP, r *http.Request) {
	context, err := m.limiter.Peek(limiterIP.String())
	if err != nil {
		log.WithError(err).WithFields(log.Fields{
			"ip_address":  limiterIP.String(),
			"request_uri": r.RequestURI,
		}).Error("Failed to get limiter context from request")
	} else {
		requestsPerformed := context.Limit - context.Remaining
		limitState := requestsPerformed * 100 / context.Limit

		operation := bucket.BuildHTTPRequestMetricOperation(r, m.statsClient.GetHTTPMetricCallback())
		// replace request method with fixed section name
		operation[0] = limiterMetric

		m.statsClient.TrackState(limiterSection, operation, int(limitState))
	}
}
