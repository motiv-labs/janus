package main

import (
	"github.com/kataras/iris"
	"gopkg.in/redis.v3"
	"net"
	"time"
	"github.com/etcinit/speedbump"
	log "github.com/Sirupsen/logrus"
	"strconv"
)

type RateLimitMiddleware struct {
	limiter *speedbump.RateLimiter
	hasher  speedbump.RateHasher
	limit   int64
}

func NewRateLimitMiddleware(client *redis.Client, hasher speedbump.RateHasher, max int64) *RateLimitMiddleware {
	limiter := speedbump.NewLimiter(client, hasher, max)
	return &RateLimitMiddleware{limiter, hasher, max}
}

//Important staff, iris middleware must implement the iris.Handler interface which is:
func (m RateLimitMiddleware) Serve(c *iris.Context) {
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

		c.Response.Header.Add("X-Rate-Limit-Limit", strconv.FormatInt(m.limit, 10))
		c.Response.Header.Add("X-Rate-Limit-Remaining", strconv.FormatInt(left, 10))
		c.Response.Header.Add("X-Rate-Limit-Reset", nextTime.String())

		c.JSON(iris.StatusTooManyRequests, map[string]string{"error": "Rate limit exceeded. Try again in " + nextTime.String()})
	} else {
		c.Next()
	}
}
