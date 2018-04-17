package retry

import (
	"net/http"
	"time"

	"github.com/Knetic/govaluate"
	"github.com/felixge/httpsnoop"
	janusErr "github.com/hellofresh/janus/pkg/errors"
	"github.com/hellofresh/janus/pkg/metrics"
	"github.com/pkg/errors"
	"github.com/rafaeljesus/retry-go"
	log "github.com/sirupsen/logrus"
)

const (
	defaultPredicate = "statusCode == 0 || statusCode >= 500"
	proxySection     = "proxy"
)

// NewRetryMiddleware creates a new retry middleware
func NewRetryMiddleware(cfg Config) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.WithFields(log.Fields{
				"attempts": cfg.Attempts,
				"backoff":  cfg.Backoff,
			}).Debug("Starting retry middleware")

			if cfg.Predicate == "" {
				cfg.Predicate = defaultPredicate
			}

			expression, err := govaluate.NewEvaluableExpression(cfg.Predicate)
			if err != nil {
				log.WithError(err).Error("could not create an expression with this predicate")
				handler.ServeHTTP(w, r)
				return
			}

			if err := retry.Do(func() error {
				m := httpsnoop.CaptureMetrics(handler, w, r)

				params := make(map[string]interface{}, 8)
				params["statusCode"] = m.Code
				params["request"] = r

				result, err := expression.Evaluate(params)
				if err != nil {
					return errors.New("cannot evaluate the expression")
				}

				if result.(bool) {
					return errors.Errorf("%s %s request failed", r.Method, r.URL)
				}

				return nil
			}, cfg.Attempts, time.Duration(cfg.Backoff)); err != nil {
				statsClient := metrics.WithContext(r.Context())
				statsClient.SetHTTPRequestSection(proxySection).TrackRequest(r, nil, false).ResetHTTPRequestSection()
				janusErr.Handler(w, errors.Wrap(err, "request failed too many times"))
			}
		})
	}
}
