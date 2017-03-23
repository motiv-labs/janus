package middleware

import (
	"net/http"
	"strings"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/hellofresh/janus/pkg/response"
)

// RateLimitLogger represents the middleware struct
type RateLimitLogger struct{}

// NewRateLimitLogger logs the IP of blocked users with rate limit
func NewRateLimitLogger() *RateLimitLogger {
	return &RateLimitLogger{}
}

// Handler is the middleware function
func (m *RateLimitLogger) Handler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			lock          sync.Mutex
			headerWritten bool
		)

		log.Debug("Starting ResponseWriterWrapper middleware")

		hooks := response.Hooks{
			WriteHeader: func(next response.WriteHeaderFunc) response.WriteHeaderFunc {
				return func(code int) {
					next(code)
					lock.Lock()
					defer lock.Unlock()
					if !headerWritten {
						if code == http.StatusTooManyRequests {
							log.WithFields(log.Fields{
								"ip_address":  realIP(r),
								"request_uri": r.RequestURI,
							}).Debug("Rate Limit exceded for this IP")
						}
						headerWritten = true
					}
				}
			},
		}

		handler.ServeHTTP(response.Wrap(w, hooks), r)
	})
}

// Request.RemoteAddress contains port, which we want to remove i.e.:
// "[::1]:58292" => "[::1]"
func ipAddrFromRemoteAddr(s string) string {
	idx := strings.LastIndex(s, ":")
	if idx == -1 {
		return s
	}
	return s[:idx]
}

// RealIP return client's real public IP address
// from http request headers.
func realIP(r *http.Request) string {
	hdr := r.Header
	hdrRealIP := hdr.Get("X-Real-Ip")
	hdrForwardedFor := hdr.Get("X-Forwarded-For")

	if len(hdrForwardedFor) == 0 && len(hdrRealIP) == 0 {
		return ipAddrFromRemoteAddr(r.RemoteAddr)
	}

	// X-Forwarded-For is potentially a list of addresses separated with ","
	for _, addr := range strings.Split(hdrForwardedFor, ",") {
		// return first non-local address
		addr = strings.TrimSpace(addr)
		if len(addr) > 0 {
			return addr
		}
	}

	return hdrRealIP
}
