package main

import (
	"errors"
	"net"
	"strconv"
	"time"

	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/etcinit/speedbump"
	"github.com/gin-gonic/gin"
)

type RateLimitMiddleware struct {
	*Middleware
	limiter *speedbump.RateLimiter
	hasher  speedbump.RateHasher
	limit   int64
}

func (m *RateLimitMiddleware) ProcessRequest(req *http.Request, c *gin.Context) (error, int) {
	log.Debug("Rate Limit middleware started")

	if !m.Spec.RateLimit.Enabled {
		log.Debug("Rate limit is not enabled")
		return nil, http.StatusOK
	}

	ip, _, _ := net.SplitHostPort(req.RemoteAddr)
	ok, err := m.limiter.Attempt(ip)
	if err != nil {
		return err, http.StatusInternalServerError
	}

	nextTime := time.Now().Add(m.hasher.Duration())
	left, err := m.limiter.Left(ip)
	if err != nil {
		return err, http.StatusInternalServerError
	}

	c.Writer.Header().Add("X-Rate-Limit-Limit", strconv.FormatInt(m.limit, 10))
	c.Writer.Header().Add("X-Rate-Limit-Remaining", strconv.FormatInt(left, 10))
	c.Writer.Header().Add("X-Rate-Limit-Reset", nextTime.String())

	if !ok {
		return errors.New("Rate limit exceeded. Try again in " + nextTime.String()), http.StatusTooManyRequests
	}

	return nil, http.StatusOK
}
