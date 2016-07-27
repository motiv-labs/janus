package main

import (
	"errors"
	"net"
	"strconv"
	"time"

	"github.com/etcinit/speedbump"
	"github.com/kataras/iris"
	"github.com/valyala/fasthttp"
)

type RateLimitMiddleware struct {
	*Middleware
	limiter *speedbump.RateLimiter
	hasher  speedbump.RateHasher
	limit   int64
}

//Important staff, iris middleware must implement the iris.Handler interface which is:
func (m RateLimitMiddleware) ProcessRequest(req fasthttp.Request, resp fasthttp.Response, c *iris.Context) (error, int) {
	m.Logger.Debug("Rate Limit middleware started")

	if !m.Spec.RateLimit.Enabled {
		m.Logger.Debug("Rate limit is not enabled")
		return nil, fasthttp.StatusOK
	}

	ip, _, _ := net.SplitHostPort(c.RemoteAddr())
	ok, err := m.limiter.Attempt(ip)
	if err != nil {
		m.Logger.Panic(err)
	}

	nextTime := time.Now().Add(m.hasher.Duration())
	left, err := m.limiter.Left(ip)
	if err != nil {
		m.Logger.Panic(err)
	}

	resp.Header.Add("X-Rate-Limit-Limit", strconv.FormatInt(m.limit, 10))
	resp.Header.Add("X-Rate-Limit-Remaining", strconv.FormatInt(left, 10))
	resp.Header.Add("X-Rate-Limit-Reset", nextTime.String())

	if !ok {
		m.Logger.Debug("Rate limit exceeded.")
		return errors.New("Rate limit exceeded. Try again in " + nextTime.String()), iris.StatusTooManyRequests
	} else {
		return nil, fasthttp.StatusOK
	}
}
