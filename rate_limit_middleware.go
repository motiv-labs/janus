package janus

import (
	"errors"
	"net"
	"strconv"
	"time"

	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/etcinit/speedbump"
	"github.com/hellofresh/janus/response"
)

// RateLimitMiddleware prevents requests to an API from exceeding a specified rate limit.
type RateLimitMiddleware struct {
	*Middleware
	limiter *speedbump.RateLimiter
	hasher  speedbump.RateHasher
	limit   int64
}

// ProcessRequest is the middleware method.
func (m *RateLimitMiddleware) Serve(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Rate Limit middleware started")

		if !m.Spec.RateLimit.Enabled {
			log.Debug("Rate limit is not enabled")
			handler.ServeHTTP(w, r)
		}

		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		ok, err := m.limiter.Attempt(ip)
		if err != nil {
			response.JSON(w, http.StatusInternalServerError, err)
			return
		}

		nextTime := time.Now().Add(m.hasher.Duration())
		left, err := m.limiter.Left(ip)
		if err != nil {
			response.JSON(w, http.StatusInternalServerError, err)
			return
		}

		w.Header().Set("X-Rate-Limit-Limit", strconv.FormatInt(m.limit, 10))
		w.Header().Set("X-Rate-Limit-Remaining", strconv.FormatInt(left, 10))
		w.Header().Set("X-Rate-Limit-Reset", nextTime.String())

		if !ok {
			response.JSON(w, http.StatusTooManyRequests, errors.New("rate limit exceeded. Try again in "+nextTime.String()))
			return
		}

		handler.ServeHTTP(w, r)
	})
}
