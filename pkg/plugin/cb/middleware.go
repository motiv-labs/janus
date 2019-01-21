package cb

import (
	"net/http"

	"github.com/Knetic/govaluate"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/felixge/httpsnoop"
	janusErr "github.com/hellofresh/janus/pkg/errors"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const (
	defaultPredicate = "statusCode == 0 || statusCode >= 500"
)

// NewCBMiddleware creates a new cb middleware
func NewCBMiddleware(cfg Config) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger := log.WithFields(log.Fields{
				"name":                    cfg.Name,
				"timeout":                 cfg.Timeout,
				"max_concurrent_requests": cfg.MaxConcurrentRequests,
				"error_percent_threshold": cfg.ErrorPercentThreshold,
			})

			logger.Debug("Starting cb middleware")
			if cfg.Predicate == "" {
				cfg.Predicate = defaultPredicate
			}

			expression, err := govaluate.NewEvaluableExpression(cfg.Predicate)
			if err != nil {
				log.WithError(err).Error("could not create an expression with this predicate")
				handler.ServeHTTP(w, r)
				return
			}

			err = hystrix.Do(cfg.Name, func() error {
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
			}, nil)

			if err != nil {
				logger.WithError(err).Error("Request failed on the cb middleware")
				janusErr.Handler(w, r, errors.Wrap(err, "request failed"))
			}
		})
	}
}
