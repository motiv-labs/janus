package janus

import (
	"errors"
	"net"
	"strconv"
	"time"

	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/etcinit/speedbump"
)

// RateLimitMiddleware prevents requests to an API from exceeding a specified rate limit.
type RateLimitMiddleware struct {
	*Middleware
	limiter *speedbump.RateLimiter
	hasher  speedbump.RateHasher
	limit   int64
}

// ProcessRequest is the middleware method.
func (m *RateLimitMiddleware) ProcessRequest(req *http.Request, rw *http.ResponseWriter) (int, error) {
	log.Debug("Rate Limit middleware started")

	if !m.Spec.RateLimit.Enabled {
		log.Debug("Rate limit is not enabled")
		return http.StatusOK, nil
	}

	ip, _, _ := net.SplitHostPort(req.RemoteAddr)
	ok, err := m.limiter.Attempt(ip)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	nextTime := time.Now().Add(m.hasher.Duration())
	left, err := m.limiter.Left(ip)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	rw.Header().Set("X-Rate-Limit-Limit", strconv.FormatInt(m.limit, 10))
	rw.Header().Set("X-Rate-Limit-Remaining", strconv.FormatInt(left, 10))
	rw.Header().Set("X-Rate-Limit-Reset", nextTime.String())

	if !ok {
		return http.StatusTooManyRequests, errors.New("Rate limit exceeded. Try again in " + nextTime.String())
	}

	return http.StatusOK, nil
}
