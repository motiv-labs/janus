package main

import (
	"github.com/kataras/iris"
	"net"
	"time"
	"github.com/etcinit/speedbump"
	log "github.com/Sirupsen/logrus"
	"strconv"
	"github.com/valyala/fasthttp"
	"errors"
)

type RateLimitMiddleware struct {
	*Middleware
	limiter *speedbump.RateLimiter
	hasher  speedbump.RateHasher
	limit   int64
}

//Important staff, iris middleware must implement the iris.Handler interface which is:
func (m RateLimitMiddleware) ProcessRequest(req fasthttp.Request, resp fasthttp.Response, c *iris.Context) (error, int) {

	if !m.Spec.RateLimit.Enabled {
		return nil, fasthttp.StatusOK
	}

	ip, _, _ := net.SplitHostPort(c.RemoteAddr())
	ok, err := m.limiter.Attempt(ip)
	if err != nil {
		log.Panic(err)
	}

	if !ok {
		nextTime := time.Now().Add(m.hasher.Duration())

		left, err := m.limiter.Left(ip)
		if err != nil {
			log.Panic(err)
		}

		resp.Header.Add("X-Rate-Limit-Limit", strconv.FormatInt(m.limit, 10))
		resp.Header.Add("X-Rate-Limit-Remaining", strconv.FormatInt(left, 10))
		resp.Header.Add("X-Rate-Limit-Reset", nextTime.String())

		return errors.New("Rate limit exceeded. Try again in " + nextTime.String()), iris.StatusTooManyRequests
	} else {
		return nil, fasthttp.StatusOK
	}
}
